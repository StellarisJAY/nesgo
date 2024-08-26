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

type SaveMetadata struct {
	Id         int64  `json:"id"`
	RoomId     int64  `json:"roomId"`
	Game       string `json:"game"`
	CreateTime int64  `json:"createTime"`
	ExitSave   bool   `json:"exitSave"`
}

type GraphicOptions struct {
	HighResOpen  bool `json:"highResOpen"`
	ReverseColor bool `json:"reverseColor"`
	Grayscale    bool `json:"grayscale"`
}

type GamingRepo interface {
	ListGames(ctx context.Context) ([]*GameMetadata, error)
	DeleteMemberConnection(ctx context.Context, roomId, userId int64, endpoint string) error
	RestartEmulator(ctx context.Context, roomId int64, game string, endpoint string) error
	ListSaves(ctx context.Context, roomId int64, page, pageSize int32) ([]*SaveMetadata, int32, error)
	SaveGame(ctx context.Context, roomId int64, endpoint string) error
	LoadSave(ctx context.Context, roomId, saveId int64, endpoint string) error
	DeleteSave(ctx context.Context, saveId int64, endpoint string) error
	GetServerICECandidate(ctx context.Context, roomId, userId int64, endpoint string) ([]string, error)
	GetGraphicOptions(ctx context.Context, roomId int64, endpoint string) (*GraphicOptions, error)
	SetGraphicOptions(ctx context.Context, roomId int64, options *GraphicOptions, endpoint string) error
	SetEmulatorSpeed(ctx context.Context, roomId int64, rate float64, endpoint string) (float64, error)
	GetEmulatorSpeed(ctx context.Context, roomId int64, endpoint string) (float64, error)
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
		UserId:     playerId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (uc *GamingUseCase) RestartEmulator(ctx context.Context, roomId, userId int64, game string) error {
	room, _ := uc.roomRepo.GetRoom(ctx, roomId)
	if room == nil {
		return v1.ErrorOperationFailed("room not found")
	}
	if room.Host != userId {
		return v1.ErrorOperationFailed("only host can restart nes")
	}
	session, _ := uc.roomRepo.GetRoomSession(ctx, roomId)
	if session == nil {
		return v1.ErrorOperationFailed("room session not found")
	}
	return uc.repo.RestartEmulator(ctx, roomId, game, session.Endpoint)
}

func (uc *GamingUseCase) ListSaves(ctx context.Context, roomId int64, page, pageSize int32) ([]*SaveMetadata, int32, error) {
	return uc.repo.ListSaves(ctx, roomId, page, pageSize)
}

func (uc *GamingUseCase) SaveGame(ctx context.Context, roomId, userId int64) error {
	room, _ := uc.roomRepo.GetRoom(ctx, roomId)
	if room == nil {
		return v1.ErrorOperationFailed("room not found")
	}
	if room.Host != userId {
		return v1.ErrorOperationFailed("only host can restart nes")
	}
	session, _ := uc.roomRepo.GetRoomSession(ctx, roomId)
	if session == nil {
		return v1.ErrorOperationFailed("room session not found")
	}
	return uc.repo.SaveGame(ctx, roomId, session.Endpoint)
}

func (uc *GamingUseCase) LoadSave(ctx context.Context, roomId, saveId int64, userId int64) error {
	room, _ := uc.roomRepo.GetRoom(ctx, roomId)
	if room == nil {
		return v1.ErrorOperationFailed("room not found")
	}
	if room.Host != userId {
		return v1.ErrorOperationFailed("only host can restart nes")
	}
	session, _ := uc.roomRepo.GetRoomSession(ctx, roomId)
	if session == nil {
		return v1.ErrorOperationFailed("room session not found")
	}
	return uc.repo.LoadSave(ctx, roomId, saveId, session.Endpoint)
}

func (uc *GamingUseCase) DeleteSave(ctx context.Context, roomId, saveId, userId int64) error {
	room, _ := uc.roomRepo.GetRoom(ctx, roomId)
	if room == nil {
		return v1.ErrorOperationFailed("room not found")
	}
	if room.Host != userId {
		return v1.ErrorOperationFailed("only host can restart nes")
	}
	session, _ := uc.roomRepo.GetRoomSession(ctx, roomId)
	if session == nil {
		return v1.ErrorOperationFailed("room session not found")
	}
	return uc.repo.DeleteSave(ctx, saveId, session.Endpoint)
}

func (uc *GamingUseCase) GetServerICECandidate(ctx context.Context, roomId, userId int64) ([]string, error) {
	session, _ := uc.roomRepo.GetRoomSession(ctx, roomId)
	if session == nil {
		return nil, v1.ErrorOperationFailed("room session not found")
	}
	candidates, err := uc.repo.GetServerICECandidate(ctx, roomId, userId, session.Endpoint)
	if err != nil {
		return nil, err
	}
	return candidates, nil
}

func (uc *GamingUseCase) GetGraphicOptions(ctx context.Context, roomId int64) (*GraphicOptions, error) {
	session, _ := uc.roomRepo.GetRoomSession(ctx, roomId)
	if session == nil {
		return nil, v1.ErrorOperationFailed("room session not found")
	}
	options, err := uc.repo.GetGraphicOptions(ctx, roomId, session.Endpoint)
	if err != nil {
		return nil, err
	}
	return options, nil
}

func (uc *GamingUseCase) SetGraphicOptions(ctx context.Context, roomId, userId int64, options *GraphicOptions) error {
	room, _ := uc.roomRepo.GetRoom(ctx, roomId)
	if room == nil {
		return v1.ErrorOperationFailed("room not found")
	}
	if room.Host != userId {
		return v1.ErrorOperationFailed("only host can restart nes")
	}
	session, _ := uc.roomRepo.GetRoomSession(ctx, roomId)
	if session == nil {
		return v1.ErrorOperationFailed("room session not found")
	}
	return uc.repo.SetGraphicOptions(ctx, roomId, options, session.Endpoint)
}

func (uc *GamingUseCase) GetEmulatorSpeed(ctx context.Context, roomId int64) (float64, error) {
	session, _ := uc.roomRepo.GetRoomSession(ctx, roomId)
	if session == nil {
		return 0, v1.ErrorOperationFailed("room session not found")
	}
	return uc.repo.GetEmulatorSpeed(ctx, roomId, session.Endpoint)
}

func (uc *GamingUseCase) SetEmulatorSpeed(ctx context.Context, roomId, userId int64, rate float64) (float64, error) {
	room, _ := uc.roomRepo.GetRoom(ctx, roomId)
	if room == nil {
		return 0, v1.ErrorOperationFailed("room not found")
	}
	if room.Host != userId {
		return 0, v1.ErrorOperationFailed("only host can set emulator speed")
	}
	session, _ := uc.roomRepo.GetRoomSession(ctx, roomId)
	if session == nil {
		return 0, v1.ErrorOperationFailed("room session not found")
	}
	return uc.repo.SetEmulatorSpeed(ctx, roomId, rate, session.Endpoint)
}
