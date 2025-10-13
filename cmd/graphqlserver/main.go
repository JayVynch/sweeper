package main

import (
	"context"
	"fmt"
	"log"

	"github.com/JayVynch/sweeper/config"
	"github.com/JayVynch/sweeper/database"
)

func main() {
	ctx := context.Background()

	conf := config.New()

	db := database.New(ctx, *conf)

	if err := db.Migrate(); err != nil {
		log.Fatalf("Could not run migration: %v", err)
	}

	fmt.Println("Migration working")
}
