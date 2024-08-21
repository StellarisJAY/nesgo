package biz

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	adminAPI "github.com/stellarisJAY/nesgo/backend/api/app/admin/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/conf"
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminClaims struct {
	Id int64
	jwt.RegisteredClaims
}

type AdminRepo interface {
	GetAdmin(ctx context.Context, id int64) (*Admin, error)
	GetAdminByName(ctx context.Context, username string) (*Admin, error)
	CreateAdmin(ctx context.Context, admin *Admin) error
}

type AdminUseCase struct {
	repo      AdminRepo
	logger    *log.Helper
	secretKey string
}

func NewAdminUseCase(repo AdminRepo, ac *conf.Auth, logger log.Logger) *AdminUseCase {
	return &AdminUseCase{
		repo:      repo,
		logger:    log.NewHelper(log.With(logger, "module", "biz/admin")),
		secretKey: ac.Secret,
	}
}

func (u *AdminUseCase) GetAdmin(ctx context.Context, id int64) (*Admin, error) {
	admin, err := u.repo.GetAdmin(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (u *AdminUseCase) CreateAdmin(ctx context.Context, admin *Admin) error {
	sum := md5.Sum([]byte(admin.Password))
	password := hex.EncodeToString(sum[:])
	admin.Password = password
	return u.repo.CreateAdmin(ctx, admin)
}

func (u *AdminUseCase) Login(ctx context.Context, name string, password string) (string, error) {
	admin, err := u.repo.GetAdminByName(ctx, name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", adminAPI.ErrorLoginFailed("admin not found")
	}
	if err != nil {
		return "", err
	}
	sum := md5.Sum([]byte(password))
	password = hex.EncodeToString(sum[:])
	if admin.Password != password {
		return "", adminAPI.ErrorLoginFailed("wrong password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AdminClaims{
		Id: admin.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "nesgo",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	signedString, err := token.SignedString([]byte(u.secretKey))
	if err != nil {
		return "", adminAPI.ErrorLoginFailed("generate token failed: %v", err)
	}
	return signedString, nil
}

func (a *AdminClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return a.ExpiresAt, nil
}

func (a *AdminClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return a.IssuedAt, nil
}

func (a *AdminClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return a.NotBefore, nil
}

func (a *AdminClaims) GetIssuer() (string, error) {
	return a.Issuer, nil
}

func (a *AdminClaims) GetSubject() (string, error) {
	return a.Subject, nil
}

func (a *AdminClaims) GetAudience() (jwt.ClaimStrings, error) {
	return a.Audience, nil
}
