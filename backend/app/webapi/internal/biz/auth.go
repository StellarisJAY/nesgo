package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/conf"
)

type AuthUseCase struct {
	ur        UserRepo
	secretKey string
	logger    *log.Helper
}

func NewAuthUseCase(ur UserRepo, c *conf.Auth, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{
		ur:        ur,
		secretKey: c.Secret,
		logger:    log.NewHelper(log.With(logger, "module", "biz/auth")),
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.ur.GetUserByName(ctx, username)
	if err != nil {
		return "", v1.ErrorLoginFailed("invalid username or password")
	}
	err = uc.ur.VerifyPassword(ctx, username, password)
	if err != nil {
		return "", v1.ErrorLoginFailed("invalid username or password")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{"userId": user.ID})
	token, err := claims.SignedString(uc.secretKey)
	if err != nil {
		return "", v1.ErrorLoginFailed("invalid username or password")
	}
	return token, nil
}

func (uc *AuthUseCase) Register(ctx context.Context, username, password string) error {
	user, err := uc.ur.GetUserByName(ctx, username)
	if user != nil {
		return v1.ErrorUsernameConflict("username already exists")
	}
	_, err = uc.ur.CreateUser(ctx, &User{
		ID:       0,
		Name:     username,
		Password: password,
	})
	if err != nil {
		return v1.ErrorRegisterFailed("create user failed", err)
	}
	return nil
}
