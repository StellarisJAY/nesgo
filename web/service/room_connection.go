package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/stellarisJAY/nesgo/web/codec"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"log"
	"net/http"
	"sync/atomic"
)

type MsgWithConnectionInfo struct {
	Message
	RTCRoomConnection
}

type Message struct {
	Type byte   `json:"type"`
	Data []byte `json:"data"`
}

const (
	MessageSDPOffer byte = iota
	MessageSDPAnswer
	MessageICECandidate
	MessageGameButtonPressed
	MessageGameButtonReleased
	MessageControlTransferred
	MessageRoomMemberChange
)

type RTCRoomConnection struct {
	MemberId     int64
	wsConn       *WebsocketConn
	rtcConn      *webrtc.PeerConnection
	videoTrack   *webrtc.TrackLocalStaticSample
	audioTrack   *webrtc.TrackLocalStaticSample
	videoEncoder codec.IVideoEncoder // 每个连接独占一个视频编码器 和 buffer
	connected    *atomic.Bool
}

type WebsocketConn struct {
	Member *room.Member
	Conn   *websocket.Conn
}

type ControlTransferredNotification struct {
	Control1 int64 `json:"control1"`
	Control2 int64 `json:"control2"`
}

type RoomMemberChangeNotification struct {
	MemberId   int64 `json:"id"`
	MemberType byte  `json:"memberType"`
	Kick       bool  `json:"kick"`
}

func (rc *RTCRoomConnection) Handle(ctx context.Context, msgChan chan MsgWithConnectionInfo) {
	wsConn := rc.wsConn.Conn
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msgType, payload, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("ws read error:", err)
			return
		}
		switch msgType {
		case websocket.TextMessage:
			msg := Message{}
			if err := json.Unmarshal(payload, &msg); err != nil {
				log.Println("invalid message error:", err)
				continue
			}
			rc.HandleMessage(msg, msgChan)
		case websocket.BinaryMessage:
		}
	}
}

func (rc *RTCRoomConnection) sendMessage(msg Message) error {
	payload, _ := json.Marshal(msg)
	return rc.wsConn.Conn.WriteMessage(websocket.TextMessage, payload)
}

func (rc *RTCRoomConnection) Close() {
	_ = rc.wsConn.Conn.Close()
	_ = rc.rtcConn.Close()
}

func (rc *RTCRoomConnection) HandleMessage(msg Message, msgChan chan MsgWithConnectionInfo) {
	switch msg.Type {
	case MessageSDPAnswer:
		sdp := webrtc.SessionDescription{}
		_ = json.Unmarshal(msg.Data, &sdp)
		if err := rc.rtcConn.SetRemoteDescription(sdp); err != nil {
			log.Println("unable to set remote description, error:", err)
			rc.Close()
		}
	case MessageICECandidate:
		candidate := webrtc.ICECandidateInit{}
		_ = json.Unmarshal(msg.Data, &candidate)
		log.Println("candidate:", candidate)
		if err := rc.rtcConn.AddICECandidate(candidate); err != nil {
			log.Println("unable to add ICE candidate, error:", err)
			rc.Close()
		}
	case MessageGameButtonPressed, MessageGameButtonReleased:
		msgChan <- MsgWithConnectionInfo{
			Message:           msg,
			RTCRoomConnection: *rc,
		}
	default:
	}
}

func (rs *RoomService) ConnectRTCRoomSession(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	v, _ := c.Get("optMember")
	member := v.(*room.Member)

	rs.m.Lock()
	var session *RTCRoomSession
	// check if room's game session is created
	if s, ok := rs.rtcSessions[roomId]; !ok {
		// Only owner can create session
		if member.MemberType != room.MemberTypeOwner {
			rs.m.Unlock()
			c.JSON(200, JSONResp{
				Status:  http.StatusForbidden,
				Message: "only owner can start game session",
			})
			return
		}

		game := c.Query("game")
		if game == "" {
			rs.m.Unlock()
			c.JSON(200, JSONResp{
				Status:  http.StatusBadRequest,
				Message: "invalid game name",
			})
			return
		}
		newSession, err := NewRTCRoomSession(game)
		if err != nil {
			panic(err)
		}
		rs.rtcSessions[roomId] = newSession
		ctx, cancelFunc := context.WithCancel(context.Background())
		newSession.cancel = cancelFunc
		go newSession.ControlLoop(ctx)
		go newSession.audioSampleListener(ctx)
		session = newSession
	} else {
		session = s
	}
	rs.m.Unlock()
	// handle room websocket conn
	conn, err := websocket.Upgrade(c.Writer, c.Request, http.Header{}, 1024, 1024)
	if err != nil {
		panic(err)
	}
	session.signalChan <- Signal{
		Type: SignalNewConnection,
		Data: &WebsocketConn{
			Member: member,
			Conn:   conn,
		},
	}
}
