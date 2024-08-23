package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	userAPI "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

type macroRepo struct {
	data   *Data
	logger *log.Helper
}

func NewMacroRepo(data *Data, logger log.Logger) biz.MacroRepo {
	return &macroRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/macro")),
	}
}

func (m *macroRepo) CreateMacro(ctx context.Context, macro *biz.Macro) error {
	_, err := m.data.uc.CreateMacro(ctx, convertMacroBizToPbCreateRequest(macro))
	return err
}

func (m *macroRepo) GetMacro(ctx context.Context, id int64) (*biz.Macro, error) {
	resp, err := m.data.uc.GetMacro(ctx, &userAPI.GetMacroRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return convertMacroPbToBiz(resp.Macro), nil
}

func (m *macroRepo) ListMacro(ctx context.Context, userId int64, page, pageSize int32) ([]*biz.Macro, int32, error) {
	resp, err := m.data.uc.ListMacro(ctx, &userAPI.ListMacroRequest{
		UserId:   userId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	res := make([]*biz.Macro, 0, len(resp.Macros))
	for _, macro := range resp.Macros {
		res = append(res, convertMacroPbToBiz(macro))
	}
	return res, resp.Total, nil
}

func (m *macroRepo) DeleteMacro(ctx context.Context, id int64) error {
	_, err := m.data.uc.DeleteMacro(ctx, &userAPI.DeleteMacroRequest{Id: id})
	return err
}

func convertMacroBizToPbCreateRequest(macro *biz.Macro) *userAPI.CreateMacroRequest {
	res := &userAPI.CreateMacroRequest{
		UserId:      macro.UserId,
		Name:        macro.Name,
		KeyboardKey: macro.KeyboardKey,
		Actions:     make([]*userAPI.MacroAction, 0, len(macro.Actions)),
	}
	for _, action := range macro.Actions {
		res.Actions = append(res.Actions, &userAPI.MacroAction{
			EmulatorKey:  action.EmulatorKey,
			ReleaseDelay: action.ReleaseDelay,
		})
	}
	return res
}

func convertMacroPbToBiz(macro *userAPI.Macro) *biz.Macro {
	res := &biz.Macro{
		Id:          macro.Id,
		UserId:      macro.UserId,
		Name:        macro.Name,
		KeyboardKey: macro.KeyboardKey,
		Actions:     make([]biz.MacroAction, 0, len(macro.Actions)),
	}
	for _, action := range macro.Actions {
		res.Actions = append(res.Actions, biz.MacroAction{
			EmulatorKey:  action.EmulatorKey,
			ReleaseDelay: action.ReleaseDelay,
		})
	}
	return res
}
