package service

import (
	"context"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
)

func (ws *WebApiService) OpenGameConnection(ctx context.Context, request *v1.OpenGameConnectionRequest) (*v1.OpenGameConnectionResponse, error) {
	sdpOffer, err := ws.gc.OpenGameConnection(ctx, request.RoomId, request.UserId)
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
	err := ws.gc.SDPAnswer(ctx, request.RoomId, request.UserId, request.SdpAnswer)
	if err != nil {
		return nil, err
	}
	return &v1.SDPAnswerResponse{
		RoomId: request.RoomId,
		UserId: request.UserId,
	}, nil
}

func (ws *WebApiService) AddICECandidate(ctx context.Context, request *v1.AddICECandidateRequest) (*v1.AddICECandidateResponse, error) {
	err := ws.gc.AddICECandidate(ctx, request.RoomId, request.UserId, request.Candidate)
	if err != nil {
		return nil, err
	}
	return &v1.AddICECandidateResponse{
		RoomId: request.RoomId,
		UserId: request.UserId,
	}, nil
}
