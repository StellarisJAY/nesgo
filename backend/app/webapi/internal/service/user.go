package service

import (
	"context"
	"github.com/stellarisjay/nesgo/backend/api/app/webapi/v1"
)

func (ws *WebApiService) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	user, err := ws.uc.GetUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserResponse{
		Code:    200,
		Message: "success",
		Data: &v1.User{
			Id:   user.ID,
			Name: user.Name,
		},
	}, nil
}
