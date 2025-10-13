//go:build integration
// +build integration

package domain

import (
	"context"
	"testing"

	"github.com/JayVynch/sweeper"
	"github.com/JayVynch/sweeper/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAuthService_Register(t *testing.T) {
	validInput := sweeper.RegisterInput{
		Username:        "doubleO7",
		Name:            "James Bond",
		Email:           "james.bond007@mi.six",
		Password:        "password",
		ConfirmPassword: "password",
	}

	t.Run("can register a user", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TearDownDB(ctx, t, db)

		res, err := authService.Register(ctx, validInput)

		require.NoError(t, err)
		require.NotEmpty(t, res.User.Id)
		require.NotEmpty(t, res.User.Email)

		require.Equal(t, validInput.Email, res.User.Email)
		require.NotEqual(t, validInput.Password, res.User.Password)
	})

	t.Run("cannot register existing username", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TearDownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)

		_, err = authService.Register(ctx, sweeper.RegisterInput{
			Username:        "doubleO7",
			Name:            "Jackson Green",
			Email:           "james.bond008@mi.six",
			Password:        "password",
			ConfirmPassword: "password",
		})

		require.ErrorIs(t, err, sweeper.ErrUserNameTaken)
	})

}
