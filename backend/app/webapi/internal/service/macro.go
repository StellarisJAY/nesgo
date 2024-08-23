package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

func (ws *WebApiService) CreateMacro(ctx context.Context, request *v1.CreateMacroRequest) (*v1.CreateMacroResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	err := ws.mu.CreateMacro(ctx, convertCreateMacroPbToBiz(request, c.UserId))
	if err != nil {
		return nil, err
	}
	return &v1.CreateMacroResponse{}, nil
}

func (ws *WebApiService) GetMacro(ctx context.Context, request *v1.GetMacroRequest) (*v1.GetMacroResponse, error) {
	macro, err := ws.mu.GetMacro(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetMacroResponse{Macro: convertMacroBizToPb(macro)}, nil
}

func (ws *WebApiService) ListMacro(ctx context.Context, request *v1.ListMacroRequest) (*v1.ListMacroResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	macros, total, err := ws.mu.ListMacro(ctx, c.UserId, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.Macro, 0, len(macros))
	for _, macro := range macros {
		result = append(result, convertMacroBizToPb(macro))
	}
	return &v1.ListMacroResponse{Macros: result, Total: total}, nil
}

func (ws *WebApiService) DeleteMacro(ctx context.Context, request *v1.DeleteMacroRequest) (*v1.DeleteMacroResponse, error) {
	err := ws.mu.DeleteMacro(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteMacroResponse{}, nil
}

func convertCreateMacroPbToBiz(req *v1.CreateMacroRequest, userId int64) *biz.Macro {
	res := &biz.Macro{
		UserId:      userId,
		Name:        req.Name,
		KeyboardKey: req.KeyboardKey,
		Actions:     make([]biz.MacroAction, 0, len(req.Actions)),
	}
	for _, action := range req.Actions {
		res.Actions = append(res.Actions, biz.MacroAction{
			EmulatorKey:  action.EmulatorKey,
			ReleaseDelay: action.ReleaseDelay,
		})
	}
	return res
}

func convertMacroBizToPb(macro *biz.Macro) *v1.Macro {
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
	return res
}
