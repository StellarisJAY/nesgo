package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

func (ws *WebApiService) OpenGameConnection(ctx context.Context, request *v1.OpenGameConnectionRequest) (*v1.OpenGameConnectionResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	sdpOffer, err := ws.gc.OpenGameConnection(ctx, request.RoomId, claims.UserId, request.Game)
	if err != nil {
		return nil, err
	}
	return &v1.OpenGameConnectionResponse{
		RoomId:   request.RoomId,
		UserId:   request.UserId,
		SdpOffer: sdpOffer,
	}, nil
}

func (ws *WebApiService) SDPAnswer(ctx context.Context, request *v1.SDPAnswerRequest) (*v1.SDPAnswerResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := ws.gc.SDPAnswer(ctx, request.RoomId, claims.UserId, request.SdpAnswer)
	if err != nil {
		return nil, err
	}
	return &v1.SDPAnswerResponse{
		RoomId: request.RoomId,
		UserId: request.UserId,
	}, nil
}

func (ws *WebApiService) AddICECandidate(ctx context.Context, request *v1.AddICECandidateRequest) (*v1.AddICECandidateResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := ws.gc.AddICECandidate(ctx, request.RoomId, claims.UserId, request.Candidate)
	if err != nil {
		return nil, err
	}
	return &v1.AddICECandidateResponse{
		RoomId: request.RoomId,
		UserId: request.UserId,
	}, nil
}

func (ws *WebApiService) ListGames(ctx context.Context, _ *v1.ListGamesRequest) (*v1.ListGamesResponse, error) {
	games, err := ws.gc.ListGames(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.GameFileMetadata, 0, len(games))
	for _, g := range games {
		result = append(result, &v1.GameFileMetadata{
			Name:      g.Name,
			Mapper:    g.Mapper,
			Mirroring: g.Mirroring,
		})
	}
	return &v1.ListGamesResponse{Games: result}, nil
}

func (ws *WebApiService) SetController(ctx context.Context, request *v1.SetControllerRequest) (*v1.SetControllerResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := ws.gc.SetController(ctx, request.RoomId, claims.UserId, request.PlayerId, request.ControllerId)
	if err != nil {
		return nil, err
	}
	return &v1.SetControllerResponse{}, nil
}

func (ws *WebApiService) RestartEmulator(ctx context.Context, request *v1.RestartEmulatorRequest) (*v1.RestartEmulatorResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := ws.gc.RestartEmulator(ctx, request.RoomId, claims.UserId, request.Game)
	if err != nil {
		return nil, err
	}
	return &v1.RestartEmulatorResponse{}, nil
}

func (ws *WebApiService) ListSaves(ctx context.Context, request *v1.ListSavesRequest) (*v1.ListSavesResponse, error) {
	saves, total, err := ws.gc.ListSaves(ctx, request.RoomId, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.Save, 0, len(saves))
	for _, save := range saves {
		result = append(result, &v1.Save{
			RoomId:     save.RoomId,
			Id:         save.Id,
			Game:       save.Game,
			CreateTime: save.CreateTime,
			ExitSave:   save.ExitSave,
		})
	}
	return &v1.ListSavesResponse{Saves: result, Total: total}, nil
}

func (ws *WebApiService) SaveGame(ctx context.Context, request *v1.SaveGameRequest) (*v1.SaveGameResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := ws.gc.SaveGame(ctx, request.RoomId, claims.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.SaveGameResponse{}, nil
}

func (ws *WebApiService) LoadSave(ctx context.Context, request *v1.LoadSaveRequest) (*v1.LoadSaveResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := ws.gc.LoadSave(ctx, request.RoomId, request.SaveId, claims.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.LoadSaveResponse{}, nil
}

func (ws *WebApiService) DeleteSave(ctx context.Context, request *v1.DeleteSaveRequest) (*v1.DeleteSaveResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	err := ws.gc.DeleteSave(ctx, request.RoomId, request.SaveId, claims.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteSaveResponse{}, nil
}

func (ws *WebApiService) GetServerICECandidate(ctx context.Context, request *v1.GetServerICECandidateRequest) (*v1.GetServerICECandidateResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	candidates, err := ws.gc.GetServerICECandidate(ctx, request.RoomId, claims.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetServerICECandidateResponse{Candidates: candidates}, nil
}

func (ws *WebApiService) GetGraphicOptions(ctx context.Context, request *v1.GetGraphicOptionsRequest) (*v1.GetGraphicOptionsResponse, error) {
	options, err := ws.gc.GetGraphicOptions(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	return &v1.GetGraphicOptionsResponse{
		HighResOpen:  options.HighResOpen,
		ReverseColor: options.ReverseColor,
		Grayscale:    options.Grayscale,
	}, nil
}

func (ws *WebApiService) SetGraphicOptions(ctx context.Context, request *v1.SetGraphicOptionsRequest) (*v1.SetGraphicOptionsResponse, error) {
	c, _ := jwt.FromContext(ctx)
	claims := c.(*biz.LoginClaims)
	opts := &biz.GraphicOptions{
		HighResOpen:  request.HighResOpen,
		ReverseColor: request.ReverseColor,
		Grayscale:    request.Grayscale,
	}
	err := ws.gc.SetGraphicOptions(ctx, request.RoomId, claims.UserId, opts)
	if err != nil {
		return nil, err
	}
	return &v1.SetGraphicOptionsResponse{
		HighResOpen:  opts.HighResOpen,
		ReverseColor: opts.ReverseColor,
		Grayscale:    opts.Grayscale,
	}, nil
}
