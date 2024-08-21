package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	userAPI "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

type userKeyboardBindingRepo struct {
	data   *Data
	logger *log.Helper
}

func NewUserKeyboardBindingRepo(data *Data, logger log.Logger) biz.UserKeyboardBindingRepo {
	return &userKeyboardBindingRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/userKeyboardBinding")),
	}
}

func (u *userKeyboardBindingRepo) CreateKeyboardBinding(ctx context.Context, ub *biz.UserKeyboardBinding) error {
	request := &userAPI.CreateUserKeyboardBindingRequest{
		Name:     ub.Name,
		UserId:   ub.UserId,
		Bindings: make([]*userAPI.KeyboardBinding, 0, len(ub.KeyboardBindings)),
	}
	for _, key := range ub.KeyboardBindings {
		request.Bindings = append(request.Bindings, &userAPI.KeyboardBinding{
			KeyboardKey: key.KeyboardKey,
			EmulatorKey: key.EmulatorKey,
		})
	}
	_, err := u.data.uc.CreateUserKeyboardBinding(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (u *userKeyboardBindingRepo) UpdateKeyboardBinding(ctx context.Context, ub *biz.UserKeyboardBinding) error {
	request := &userAPI.UpdateUserKeyboardBindingRequest{
		Id:       ub.Id,
		Name:     ub.Name,
		UserId:   ub.UserId,
		Bindings: make([]*userAPI.KeyboardBinding, 0, len(ub.KeyboardBindings)),
	}
	for _, key := range ub.KeyboardBindings {
		request.Bindings = append(request.Bindings, &userAPI.KeyboardBinding{
			KeyboardKey: key.KeyboardKey,
			EmulatorKey: key.EmulatorKey,
		})
	}
	_, err := u.data.uc.UpdateUserKeyboardBinding(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (u *userKeyboardBindingRepo) DeleteKeyboardBinding(ctx context.Context, id int64) error {
	_, err := u.data.uc.DeleteUserKeyboardBinding(ctx, &userAPI.DeleteUserKeyboardBindingRequest{
		Id: id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userKeyboardBindingRepo) GetKeyboardBinding(ctx context.Context, id int64) (*biz.UserKeyboardBinding, error) {
	response, err := u.data.uc.GetUserKeyboardBinding(ctx, &userAPI.GetUserKeyboardBindingRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return convertPbToBiz(response.Binding), nil
}

func (u *userKeyboardBindingRepo) ListUserKeyboardBinding(ctx context.Context, userId int64, page, pageSize int32) ([]*biz.UserKeyboardBinding, int32, error) {
	response, err := u.data.uc.ListUserKeyboardBinding(ctx, &userAPI.ListUserKeyboardBindingRequest{
		UserId:   userId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	res := make([]*biz.UserKeyboardBinding, 0, len(response.Bindings))
	for _, item := range response.Bindings {
		res = append(res, convertPbToBiz(item))
	}
	return res, response.Total, nil
}

func convertBizToPb(ub *biz.UserKeyboardBinding) *userAPI.UserKeyboardBinding {
	userBinding := &userAPI.UserKeyboardBinding{
		Id:       ub.Id,
		Name:     ub.Name,
		UserId:   ub.UserId,
		Bindings: make([]*userAPI.KeyboardBinding, 0, len(ub.KeyboardBindings)),
	}
	for _, key := range ub.KeyboardBindings {
		userBinding.Bindings = append(userBinding.Bindings, &userAPI.KeyboardBinding{
			KeyboardKey: key.KeyboardKey,
			EmulatorKey: key.EmulatorKey,
		})
	}
	return userBinding
}

func convertPbToBiz(ub *userAPI.UserKeyboardBinding) *biz.UserKeyboardBinding {
	userBinding := &biz.UserKeyboardBinding{
		Id:               ub.Id,
		Name:             ub.Name,
		UserId:           ub.UserId,
		KeyboardBindings: make([]*biz.KeyboardBinding, 0, len(ub.Bindings)),
	}
	for _, key := range ub.Bindings {
		userBinding.KeyboardBindings = append(userBinding.KeyboardBindings, &biz.KeyboardBinding{
			KeyboardKey: key.KeyboardKey,
			EmulatorKey: key.EmulatorKey,
		})
	}
	return userBinding
}
