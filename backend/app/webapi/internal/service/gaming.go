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
