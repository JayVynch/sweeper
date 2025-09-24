package sweeper

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var (
	UsernameMinLen = 2
	PasswordMinLen = 6
	emailRegexp    = regexp.MustCompile(`^(?:(?:[a-zA-Z0-9._%+-]+)|(?:"(?:[^"\\]|\\.)+"))@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
}

type AuthResponse struct {
	AccessToken string
	User        User
}

type RegisterInput struct {
	Name            string
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

func (in *RegisterInput) Sanitize() {
	in.Email = strings.TrimSpace(in.Email)
	in.Email = strings.ToLower(in.Email)

	in.Username = strings.TrimSpace(in.Username)
}

func (in RegisterInput) validate() error {
	if len(in.Username) < UsernameMinLen {
		return fmt.Errorf("%w: username not long enough, (%d) characters at least", ErrValidation, UsernameMinLen)
	}

	if !emailRegexp.MatchString(in.Email) {
		return fmt.Errorf("%w:  not a valid email structure", ErrValidation)
	}

	if len(in.Password) < PasswordMinLen {
		return fmt.Errorf("%w: password not long enough, (%d) characters at least", ErrValidation, PasswordMinLen)
	}

	if in.Password != in.ConfirmPassword {
		return fmt.Errorf("%w: password does not match", ErrValidation)
	}
	return nil
}
