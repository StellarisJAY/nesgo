package service

import (
	"context"
	v1 "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/user/internal/biz"
)

func (u *UserService) CreateUserKeyboardBinding(ctx context.Context, request *v1.CreateUserKeyboardBindingRequest) (*v1.CreateUserKeyboardBindingResponse, error) {
	ub := &biz.UserKeyboardBinding{
		Id:               0,
		Name:             request.Name,
		UserId:           request.UserId,
		KeyboardBindings: make([]*biz.KeyboardBinding, 0, len(request.Bindings)),
	}
	for _, binding := range request.Bindings {
		ub.KeyboardBindings = append(ub.KeyboardBindings, &biz.KeyboardBinding{
			KeyboardKey: binding.KeyboardKey,
			EmulatorKey: binding.EmulatorKey,
		})
	}
	err := u.kb.CreateKeyboardBinding(ctx, ub)
	if err != nil {
		return nil, err
	}
	return &v1.CreateUserKeyboardBindingResponse{}, nil
}

func (u *UserService) ListUserKeyboardBinding(ctx context.Context, request *v1.ListUserKeyboardBindingRequest) (*v1.ListUserKeyboardBindingResponse, error) {
	ubs, total, err := u.kb.ListUserKeyboardBinding(ctx, request.UserId, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	res := make([]*v1.UserKeyboardBinding, 0, len(ubs))
	for _, ub := range ubs {
		userBinding := convertBizUserKeyboardBindingToProtoUserKeyboardBinding(ub)
		res = append(res, userBinding)
	}
	return &v1.ListUserKeyboardBindingResponse{
		Total:    total,
		Bindings: res,
	}, nil
}

func (u *UserService) GetUserKeyboardBinding(ctx context.Context, request *v1.GetUserKeyboardBindingRequest) (*v1.GetUserKeyboardBindingResponse, error) {
	ub, err := u.kb.GetKeyboardBinding(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserKeyboardBindingResponse{
		Binding: convertBizUserKeyboardBindingToProtoUserKeyboardBinding(ub),
	}, nil
}

func (u *UserService) UpdateUserKeyboardBinding(ctx context.Context, request *v1.UpdateUserKeyboardBindingRequest) (*v1.UpdateUserKeyboardBindingResponse, error) {
	userBinding := &biz.UserKeyboardBinding{
		Id:               request.Id,
		Name:             request.Name,
		UserId:           request.UserId,
		KeyboardBindings: make([]*biz.KeyboardBinding, 0, len(request.Bindings)),
	}
	for _, binding := range request.Bindings {
		userBinding.KeyboardBindings = append(userBinding.KeyboardBindings, &biz.KeyboardBinding{
			KeyboardKey: binding.KeyboardKey,
			EmulatorKey: binding.EmulatorKey,
		})
	}
	err := u.kb.UpdateKeyboardBinding(ctx, userBinding)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateUserKeyboardBindingResponse{}, nil
}

func (u *UserService) DeleteUserKeyboardBinding(ctx context.Context, request *v1.DeleteUserKeyboardBindingRequest) (*v1.DeleteUserKeyboardBindingResponse, error) {
	err := u.kb.DeleteKeyboardBinding(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteUserKeyboardBindingResponse{}, nil
}

func convertBizUserKeyboardBindingToProtoUserKeyboardBinding(ub *biz.UserKeyboardBinding) *v1.UserKeyboardBinding {
	userBinding := &v1.UserKeyboardBinding{
		Id:       ub.Id,
		Name:     ub.Name,
		UserId:   ub.UserId,
		Bindings: make([]*v1.KeyboardBinding, 0, len(ub.KeyboardBindings)),
	}
	for _, binding := range ub.KeyboardBindings {
		userBinding.Bindings = append(userBinding.Bindings, &v1.KeyboardBinding{
			KeyboardKey: binding.KeyboardKey,
			EmulatorKey: binding.EmulatorKey,
		})
	}
	return userBinding
}
