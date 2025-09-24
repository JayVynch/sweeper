package sweeper

import "time"

type UserRepo interface {
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
