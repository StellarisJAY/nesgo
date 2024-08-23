package service

import (
	"context"
	v1 "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/user/internal/biz"
)

func (u *UserService) CreateMacro(ctx context.Context, request *v1.CreateMacroRequest) (*v1.CreateMacroResponse, error) {
	macro := &biz.Macro{
		UserId:      request.UserId,
		Name:        request.Name,
		KeyboardKey: request.KeyboardKey,
		Actions:     make([]biz.MacroAction, 0, len(request.Actions)),
	}
	for _, action := range request.Actions {
		macro.Actions = append(macro.Actions, biz.MacroAction{
			EmulatorKey:  action.EmulatorKey,
			ReleaseDelay: action.ReleaseDelay,
		})
	}
	err := u.mu.CreateMacro(ctx, macro)
	if err != nil {
		return nil, err
	}
	return &v1.CreateMacroResponse{}, nil
}

func (u *UserService) GetMacro(ctx context.Context, request *v1.GetMacroRequest) (*v1.GetMacroResponse, error) {
	macro, err := u.mu.GetMacro(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	res := &v1.Macro{
		Id:          macro.Id,
		UserId:      macro.UserId,
		Name:        macro.Name,
		KeyboardKey: macro.KeyboardKey,
		Actions:     make([]*v1.MacroAction, 0, len(macro.Actions)),
	}
	for _, action := range macro.Actions {
		res.Actions = append(res.Actions, &v1.MacroAction{
			EmulatorKey:  action.EmulatorKey,
			ReleaseDelay: action.ReleaseDelay,
		})
	}
	return &v1.GetMacroResponse{Macro: res}, nil
}

func (u *UserService) ListMacro(ctx context.Context, request *v1.ListMacroRequest) (*v1.ListMacroResponse, error) {
	macros, total, err := u.mu.ListMacro(ctx, request.UserId, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	res := &v1.ListMacroResponse{
		Total:  total,
		Macros: make([]*v1.Macro, 0, len(macros)),
	}
	for _, macro := range macros {
		res.Macros = append(res.Macros, &v1.Macro{
			Id:          macro.Id,
			UserId:      macro.UserId,
			Name:        macro.Name,
			KeyboardKey: macro.KeyboardKey,
			Actions:     make([]*v1.MacroAction, 0, len(macro.Actions)),
		})
		for _, action := range macro.Actions {
			res.Macros[len(res.Macros)-1].Actions = append(res.Macros[len(res.Macros)-1].Actions, &v1.MacroAction{
				EmulatorKey:  action.EmulatorKey,
				ReleaseDelay: action.ReleaseDelay,
			})
		}
	}
	return res, nil
}

func (u *UserService) DeleteMacro(ctx context.Context, request *v1.DeleteMacroRequest) (*v1.DeleteMacroResponse, error) {
	err := u.mu.DeleteMacro(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteMacroResponse{}, nil
}
