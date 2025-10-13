//go:build integration
// +build integration

package domain

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/JayVynch/sweeper/config"
	"github.com/JayVynch/sweeper/database"
)

var (
	conf        *config.Config
	db          *database.DB
	authService AuthService
	userRepo    database.UserRepo
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	config.LoadEnv(".env.test")
	// passwordCost := bcrypt.MinCost
	conf = config.New()

	db = database.New(ctx, *conf)

	if err := db.Drop(); err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	userRepo = *database.NewUserRepo(db)

	authService = *NewAuthService(&userRepo)

	os.Exit(m.Run())
}
