package sweeper

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNameTaken = errors.New("username Taken")
	ErrEmailTaken    = errors.New("email Taken")
)

type UserRepo interface {
	Create(cxt context.Context, user User) (User, error)
	GetByUsername(cxt context.Context, username string)
	GetByEmail(cxt context.Context, email string)
}

type User struct {
	Id        string
	Username  string
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
