package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
	gamingAPI "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
)

type GamingUseCase struct {
	repo     GamingRepo
	roomRepo RoomRepo
	logger   *log.Helper
}

type GameMetadata struct {
	Name      string `json:"name"`
	Mapper    string `json:"mapper"`
	Mirroring string `json:"mirroring"`
}

type GamingRepo interface {
	ListGames(ctx context.Context) ([]*GameMetadata, error)
}

func NewGamingUseCase(roomRepo RoomRepo, gamingRepo GamingRepo, logger log.Logger) *GamingUseCase {
	return &GamingUseCase{
		repo:     gamingRepo,
		roomRepo: roomRepo,
		logger:   log.NewHelper(log.With(logger, "module", "biz/gaming")),
	}
}

func (uc *GamingUseCase) OpenGameConnection(ctx context.Context, roomId int64, userId int64, game string) (string, error) {
	// 调用房间服务，获取房间模拟器会话地址
	session, err := uc.roomRepo.GetCreateRoomSession(ctx, roomId, userId, game)
	if err != nil {
		return "", err
	}
	endpoint := session.Endpoint
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(endpoint))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	gamingCli := gamingAPI.NewGamingClient(conn)
	// 调用模拟器服务，创建webrtc连接，获取sdpOffer
	response, err := gamingCli.OpenGameConnection(ctx, &gamingAPI.OpenGameConnectionRequest{RoomId: roomId, UserId: userId})
	if err != nil {
		return "", err
	}
	return response.SdpOffer, nil
}

func (uc *GamingUseCase) SDPAnswer(ctx context.Context, roomId, userId int64, sdpAnswer string) error {
	session, err := uc.roomRepo.GetRoomSession(ctx, roomId)
	if err != nil {
		return err
	}
	endpoint := session.Endpoint
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(endpoint))
	if err != nil {
		return err
	}
	defer conn.Close()
	gamingCli := gamingAPI.NewGamingClient(conn)
	_, err = gamingCli.SDPAnswer(ctx, &gamingAPI.SDPAnswerRequest{RoomId: roomId, UserId: userId, SdpAnswer: sdpAnswer})
	return err
}

func (uc *GamingUseCase) AddICECandidate(ctx context.Context, roomId, userId int64, candidate string) error {
	session, err := uc.roomRepo.GetRoomSession(ctx, roomId)
	if err != nil {
		return err
	}
	endpoint := session.Endpoint
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(endpoint))
	if err != nil {
		return err
	}
	defer conn.Close()
	gamingCli := gamingAPI.NewGamingClient(conn)
	_, err = gamingCli.ICECandidate(ctx, &gamingAPI.ICECandidateRequest{RoomId: roomId, UserId: userId, Candidate: candidate})
	return err
}

func (uc *GamingUseCase) ListGames(ctx context.Context) ([]*GameMetadata, error) {
	return uc.repo.ListGames(ctx)
}

func (uc *GamingUseCase) SetController(ctx context.Context, roomId, userId, playerId int64, controller int32) error {
	room, err := uc.roomRepo.GetRoom(ctx, roomId)
	if err != nil {
		return err
	}
	if room.Host != userId {
		return v1.ErrorOperationFailed("only host can set controller")
	}
	member, err := uc.roomRepo.GetRoomMember(ctx, roomId, playerId)
	if err != nil {
		return err
	}
	if member.Role == roomAPI.RoomRole_Observer {
		return v1.ErrorOperationFailed("can't set controller for observer")
	}
	session, err := uc.roomRepo.GetRoomSession(ctx, roomId)
	if err != nil {
		return err
	}
	endpoint := session.Endpoint
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(endpoint))
	if err != nil {
		return v1.ErrorOperationFailed("can't connect room session: %v", err)
	}
	defer conn.Close()
	client := gamingAPI.NewGamingClient(conn)
	_, err = client.SetController(ctx, &gamingAPI.SetControllerRequest{
		Controller: controller,
		RoomId:     roomId,
		UserId:     userId,
	})
	if err != nil {
		return err
	}
	return nil
}
