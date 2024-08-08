package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/conf"
	"time"
)

type AuthUseCase struct {
	ur        UserRepo
	secretKey string
	logger    *log.Helper
}

type LoginClaims struct {
	jwt.RegisteredClaims
	UserId   int64
	Username string
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
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &LoginClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "nesgo",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId:   user.ID,
		Username: username,
	})
	token, err := claims.SignedString([]byte(uc.secretKey))
	if err != nil {
		return "", v1.ErrorLoginFailed("generate token error: %v", err)
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

func (l *LoginClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return l.ExpiresAt, nil
}

func (l *LoginClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return l.IssuedAt, nil
}

func (l *LoginClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return l.NotBefore, nil
}

func (l *LoginClaims) GetIssuer() (string, error) {
	return l.Issuer, nil
}

func (l *LoginClaims) GetSubject() (string, error) {
	return l.Subject, nil
}

func (l *LoginClaims) GetAudience() (jwt.ClaimStrings, error) {
	return l.Audience, nil
}
