package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewUserUseCase,
	NewAuthUseCase,
	NewRoomUseCase,
	NewGamingUseCase,
	NewUserKeyboardBindingUseCase,
	NewMacroUseCase,
)
