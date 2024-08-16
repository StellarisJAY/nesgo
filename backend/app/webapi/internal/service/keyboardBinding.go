package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

func (ws *WebApiService) CreateUserKeyboardBinding(ctx context.Context, request *v1.CreateUserKeyboardBindingRequest) (*v1.CreateUserKeyboardBindingResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	ub := &biz.UserKeyboardBinding{
		Name:             request.Name,
		UserId:           c.UserId,
		KeyboardBindings: make([]*biz.KeyboardBinding, 0, len(request.Bindings)),
	}
	for _, key := range request.Bindings {
		ub.KeyboardBindings = append(ub.KeyboardBindings, &biz.KeyboardBinding{
			KeyboardKey: key.KeyboardKey,
			EmulatorKey: key.EmulatorKey,
		})
	}
	err := ws.uk.Create(ctx, ub)
	if err != nil {
		return nil, err
	}
	return &v1.CreateUserKeyboardBindingResponse{}, nil
}

func (ws *WebApiService) ListUserKeyboardBinding(ctx context.Context, request *v1.ListUserKeyboardBindingRequest) (*v1.ListUserKeyboardBindingResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	ubs, total, err := ws.uk.List(ctx, c.UserId, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.UserKeyboardBinding, 0, len(ubs))
	for _, ub := range ubs {
		result = append(result, convertBizUserKeyboardBindingToPbUserKeyboardBinding(ub))
	}
	return &v1.ListUserKeyboardBindingResponse{
		Total:    total,
		Bindings: result,
	}, nil
}

func (ws *WebApiService) GetUserKeyboardBinding(ctx context.Context, request *v1.GetUserKeyboardBindingRequest) (*v1.GetUserKeyboardBindingResponse, error) {
	ub, err := ws.uk.Get(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserKeyboardBindingResponse{
		Binding: convertBizUserKeyboardBindingToPbUserKeyboardBinding(ub),
	}, nil
}

func (ws *WebApiService) DeleteUserKeyboardBinding(ctx context.Context, request *v1.DeleteUserKeyboardBindingRequest) (*v1.DeleteUserKeyboardBindingResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	err := ws.uk.Delete(ctx, request.Id, c.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteUserKeyboardBindingResponse{}, nil
}

func (ws *WebApiService) UpdateUserKeyboardBinding(ctx context.Context, request *v1.UpdateUserKeyboardBindingRequest) (*v1.UpdateUserKeyboardBindingResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	ub := &biz.UserKeyboardBinding{
		Name:             request.Name,
		UserId:           c.UserId,
		KeyboardBindings: make([]*biz.KeyboardBinding, 0, len(request.Bindings)),
	}
	for _, key := range request.Bindings {
		ub.KeyboardBindings = append(ub.KeyboardBindings, &biz.KeyboardBinding{
			KeyboardKey: key.KeyboardKey,
			EmulatorKey: key.EmulatorKey,
		})
	}
	err := ws.uk.Update(ctx, ub)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateUserKeyboardBindingResponse{}, nil
}

func convertBizUserKeyboardBindingToPbUserKeyboardBinding(ub *biz.UserKeyboardBinding) *v1.UserKeyboardBinding {
	userBindings := &v1.UserKeyboardBinding{
		Id:       ub.Id,
		Name:     ub.Name,
		UserId:   ub.UserId,
		Bindings: make([]*v1.KeyboardBinding, 0, len(ub.KeyboardBindings)),
	}
	for _, item := range ub.KeyboardBindings {
		userBindings.Bindings = append(userBindings.Bindings, &v1.KeyboardBinding{
			KeyboardKey:           item.KeyboardKey,
			EmulatorKey:           item.EmulatorKey,
			KeyboardKeyTranslated: biz.TranslateKeyboardKey(item.KeyboardKey),
			EmulatorKeyTranslated: biz.TranslateEmulatorKey(item.EmulatorKey),
		})
	}
	return userBindings
}
