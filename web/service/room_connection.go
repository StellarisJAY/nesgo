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

type MessageWithConnInfo struct {
	Message
	RoomConn
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
)

type RoomConn struct {
	MemberId     int64
	wsConn       *WebsocketConn
	rtcConn      *webrtc.PeerConnection
	videoTrack   *webrtc.TrackLocalStaticSample
	audioTrack   *webrtc.TrackLocalStaticSample
	videoEncoder codec.IVideoEncoder // 每个连接独占一个视频编码器 和 buffer
	connected    *atomic.Bool
	dataChannel  *webrtc.DataChannel
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

func (rc *RoomConn) Handle(ctx context.Context) {
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
			rc.HandleMessage(msg)
		default:
		}
	}
}

func (rc *RoomConn) sendMessage(msg Message) error {
	payload, _ := json.Marshal(msg)
	return rc.wsConn.Conn.WriteMessage(websocket.TextMessage, payload)
}

func (rc *RoomConn) Close() {
	_ = rc.wsConn.Conn.Close()
	_ = rc.rtcConn.Close()
}

func (rc *RoomConn) HandleMessage(msg Message) {
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
	default:
	}
}

func (rc *RoomConn) OnDataChannelMessage(rawMsg webrtc.DataChannelMessage, msgChan chan MessageWithConnInfo) {
	var message Message
	if err := json.Unmarshal(rawMsg.Data, &message); err != nil {
		log.Println("invalid data channel message, error:", err)
		return
	}
	if message.Type == MessageGameButtonPressed || message.Type == MessageGameButtonReleased {
		msgChan <- MessageWithConnInfo{
			Message:  message,
			RoomConn: *rc,
		}
	}
}

func (rc *RoomConn) onPeerConnectionState(state webrtc.PeerConnectionState, signalChan chan Signal) {
	log.Println("peer conn state:", state)
	switch state {
	case webrtc.PeerConnectionStateConnected:
		rc.connected.Store(true)
		signalChan <- Signal{SignalPeerConnected, rc}
	case webrtc.PeerConnectionStateDisconnected:
		signalChan <- Signal{SignalPeerDisconnected, rc}
		rc.connected.Store(false)
		rc.Close()
	default:
	}
}

func (rc *RoomConn) onICEStateChange(state webrtc.ICEConnectionState) {
	log.Println("ice conn state:", state)
}

func (rs *RoomService) ConnectRTCRoomSession(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	v, _ := c.Get("optMember")
	member := v.(*room.Member)

	rs.m.Lock()
	session, ok := rs.rtcSessions[roomId]
	rs.m.Unlock()
	if !ok {
		// Only owner can create session
		if member.MemberType != room.MemberTypeOwner {
			c.JSON(200, JSONResp{Status: http.StatusForbidden, Message: "only owner can start game session"})
			return
		}
		game := c.Query("game")
		if game == "" {
			c.JSON(200, JSONResp{Status: http.StatusBadRequest, Message: "invalid game name"})
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
	}
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
