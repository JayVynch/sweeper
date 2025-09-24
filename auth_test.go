package sweeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInputSanitize(t *testing.T) {
	input := RegisterInput{
		Name:            "Bob Sponge",
		Username:        " Bob ",
		Email:           " BOB@email.com ",
		Password:        "password",
		ConfirmPassword: "password",
	}
	want := RegisterInput{
		Name:            "Bob Sponge",
		Username:        "Bob",
		Email:           "bob@email.com",
		Password:        "password",
		ConfirmPassword: "password",
	}

	input.Sanitize()
	require.Equal(t, want, input)
}

func TestRegisterInputValidation(t *testing.T) {
	testCases := []struct {
		name  string
		input RegisterInput
		err   error
	}{
		{
			name: "valid",
			input: RegisterInput{
				Name:            "Bob Sponge",
				Username:        "Bob",
				Email:           "bob@email.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: RegisterInput{
				Name:            "Bob Sponge",
				Username:        "Bob",
				Email:           "bob@email",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		},
		{
			name: "invalid username len",
			input: RegisterInput{
				Name:            "Bob Sponge",
				Username:        "B",
				Email:           "bob@email.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		},
		{
			name: "invalid password len",
			input: RegisterInput{
				Name:            "Bob Sponge",
				Username:        "Bob",
				Email:           "bob@email.com",
				Password:        "pass",
				ConfirmPassword: "pass",
			},
			err: ErrValidation,
		},
		{
			name: "invalid password match",
			input: RegisterInput{
				Name:            "Bob Sponge",
				Username:        "Bob",
				Email:           "bob@email.com",
				Password:        "password",
				ConfirmPassword: "passwords",
			},
			err: ErrValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.validate()

			if err != nil {
				// we want error
				require.ErrorIs(t, err, tc.err)
			} else {
				// we dont want an error
				require.NoError(t, err)
			}
		})
	}
}
