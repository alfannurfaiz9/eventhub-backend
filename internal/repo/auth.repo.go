package repo

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthRepo struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (a *AuthRepo) FindUser(ctx context.Context, email string) (model.User, error) {
	sql := "SELECT id, full_name, img_url, email, password, role FROM users WHERE email=$1"
	args := []any{email}

	var data model.User
	if err := a.db.QueryRow(ctx, sql, args...).Scan(
		&data.Id,
		&data.FullName,
		&data.ImgUrl,
		&data.Email,
		&data.Password,
		&data.Role); err != nil {
		return model.User{}, err
	}

	return data, nil
}

func (a *AuthRepo) Register(ctx context.Context, body model.User) error {
	sql := "INSERT INTO users(full_name, email, password) VALUES($1, $2, $3)"
	args := []any{body.FullName, body.Email, body.Password}

	cmd, err := a.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}

func (a *AuthRepo) FindUserByEmail(ctx context.Context, body model.User) (string, error) {
	sql := `
	SELECT email
	FROM users
	WHERE email = $1`
	args := []any{body.Email}

	var data model.User
	if err := a.db.QueryRow(ctx, sql, args...).Scan(
		&data.Email,
	); err != nil {
		return "", err
	}

	return data.Email, nil
}

func (a *AuthRepo) ForgotPassword(ctx context.Context, body model.User) error {
	sql := `
	UPDATE users
	SET password = $1
	WHERE email = $2`

	args := []any{body.Password, body.Email}

	cmd, err := a.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
