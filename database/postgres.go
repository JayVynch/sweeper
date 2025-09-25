package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func new(ctx context.Context) *DB {
	dbConfig, err := pgxpool.ParseConfig()

	if err != nil {
		log.Fatalf("Cannot parse database config %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Cannot connect database %v", err)
	}

	db := &DB{Pool: pool}

	db.Ping(ctx)
	return db
}

func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatal("cannot ping database %v", err)
	}
}

func (db *DB) Open(cxt context.Context) {
	if err := db.Pool.Ping(cxt); err != nil {
		log.Fatalf("Could not ping Postgres: %v", err)
	}

	log.Println("Postgres pinged")
}

func (db *DB) Close() {
	db.Pool.Close()
}
