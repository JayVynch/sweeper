package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/JayVynch/sweeper"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo sweeper.UserRepo
}

func NewAuthService(ur sweeper.UserRepo) *AuthService {
	return &AuthService{
		UserRepo: ur,
	}
}

func (as *AuthService) Register(ctx context.Context, input sweeper.RegisterInput) (sweeper.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return sweeper.AuthResponse{}, err
	}

	// check if username is available
	if _, err := as.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, sweeper.ErrorNotFound) {
		return sweeper.AuthResponse{}, sweeper.ErrUserNameTaken
	}

	// check if email is available
	if _, err := as.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, sweeper.ErrorNotFound) {
		return sweeper.AuthResponse{}, sweeper.ErrEmailTaken
	}

	newUser := sweeper.User{
		Email:    input.Email,
		Name:     input.Name,
		Username: input.Username,
	}

	// hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return sweeper.AuthResponse{}, fmt.Errorf("error during hashing password: %w", err)
	}

	newUser.Password = string(hashedPwd)

	// create user
	newUser, err = as.UserRepo.Create(ctx, newUser)
	if err != nil {
		return sweeper.AuthResponse{}, fmt.Errorf("error creating user: %w", err)
	}
	// return AccessToken and user
	return sweeper.AuthResponse{
		AccessToken: "SweeperToken",
		User:        newUser,
	}, nil
}

func (as *AuthService) Login(ctx context.Context, input sweeper.LoginInput) (sweeper.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return sweeper.AuthResponse{}, err
	}

	user, err := as.UserRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, sweeper.ErrorNotFound):
			return sweeper.AuthResponse{}, sweeper.ErrBadCredentials
		default:
			return sweeper.AuthResponse{}, err
		}

	}

	// check password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return sweeper.AuthResponse{}, sweeper.ErrBadCredentials
	}

	// return AccessToken and user
	return sweeper.AuthResponse{
		AccessToken: "SweeperToken",
		User:        user,
	}, nil
}
