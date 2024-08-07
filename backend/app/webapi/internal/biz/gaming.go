package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gamingAPI "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
)

type GamingUseCase struct {
	repo     GamingRepo
	roomRepo RoomRepo
	logger   *log.Helper
}

type GamingRepo interface {
}

func NewGamingUseCase(roomRepo RoomRepo, gamingRepo GamingRepo, logger log.Logger) *GamingUseCase {
	return &GamingUseCase{
		repo:     gamingRepo,
		roomRepo: roomRepo,
		logger:   log.NewHelper(log.With(logger, "module", "biz/gaming")),
	}
}

func (uc *GamingUseCase) OpenGameConnection(ctx context.Context, roomId int64, userId int64) (string, error) {
	// 调用房间服务，获取房间模拟器会话地址
	session, err := uc.roomRepo.GetRoomSession(ctx, roomId, userId)
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
	session, err := uc.roomRepo.GetRoomSession(ctx, roomId, userId)
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
	session, err := uc.roomRepo.GetRoomSession(ctx, roomId, userId)
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
