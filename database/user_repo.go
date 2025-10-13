package database

import (
	"context"
	"fmt"

	"github.com/JayVynch/sweeper"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	Db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{
		Db: db,
	}
}

func (u *UserRepo) Create(ctx context.Context, user sweeper.User) (sweeper.User, error) {
	tx, err := u.Db.Pool.Begin(ctx)
	if err != nil {
		return sweeper.User{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	user, err = createUser(ctx, tx, user)
	if err != nil {
		return sweeper.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return sweeper.User{}, fmt.Errorf("error while commiting: %v", err)
	}

	return user, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user sweeper.User) (sweeper.User, error) {
	query := `INSERT INTO users (name,username,email,password) VALUES ($1,$2,$3,$4) RETURNING *;`

	ur := sweeper.User{}

	if err := pgxscan.Get(ctx, tx, &ur, query, user.Name, user.Username, user.Email, user.Password); err != nil {
		return sweeper.User{}, fmt.Errorf("error inserting fields: %v", err)
	}

	return ur, nil
}

func (u *UserRepo) GetByUsername(ctx context.Context, username string) (sweeper.User, error) {
	query := `SELECT * FROM USERS WHERE username = $1 LIMIT 1;`

	user := sweeper.User{}

	if err := pgxscan.Get(ctx, u.Db.Pool, &user, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return sweeper.User{}, sweeper.ErrorNotFound
		}

		return sweeper.User{}, fmt.Errorf("error select: %v", err)
	}

	return user, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (sweeper.User, error) {
	query := `SELECT * FROM USERS WHERE email = $1 LIMIT 1;`

	user := sweeper.User{}

	if err := pgxscan.Get(ctx, u.Db.Pool, &user, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return sweeper.User{}, sweeper.ErrorNotFound
		}

		return sweeper.User{}, fmt.Errorf("error select: %v", err)
	}

	return user, nil
}
