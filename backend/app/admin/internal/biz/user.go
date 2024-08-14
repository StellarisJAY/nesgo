package biz

import "context"

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UserRepo interface {
	GetUser(ctx context.Context, id int64) (*User, error)
}
