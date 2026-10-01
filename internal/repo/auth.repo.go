package repo

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (a *AuthRepo) FindUser(ctx context.Context, email string) (model.User, error) {
	sql := "SELECT id, email,password, role FROM users WHERE email=$1"
	args := []any{email}

	var data model.User
	if err := a.db.QueryRow(ctx, sql, args...).Scan(&data.Id, &data.Email, &data.Password, &data.Role); err != nil {
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

func (a *AuthRepo) GetAdminDashboard(ctx context.Context) (model.AdminDashboard, error) {
	sql := `
	SELECT total_users, total_events, total_communities, CONCAT((total_attendee * 100) / total_capacaity, '%')
	FROM 
		(SELECT COUNT(users.id) AS "total_users" FROM users), 
		(SELECT COUNT(events.id) AS "total_events", SUM(events.capacity) AS "total_capacaity" FROM events),
		(SELECT COUNT(communities.id) AS "total_communities" FROM communities),
		(SELECT COUNT(user_event.user_id) AS "total_attendee"  FROM user_event)`

	var data model.AdminDashboard
	err := a.db.QueryRow(ctx, sql).Scan(
		&data.TotalUser,
		&data.TotalEvent,
		&data.TotalCommunity,
		&data.AvgFillRate)

	if err != nil {
		return model.AdminDashboard{}, err
	}

	return data, nil
}

func (a *AuthRepo) GetAllUser(ctx context.Context) ([]model.User, error) {
	sql := `SELECT full_name, email, role, status, created_at FROM users`

	rows, err := a.db.Query(ctx, sql)

	if err != nil {
		return nil, err
	}

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.FullName,
			&user.Email,
			&user.Role,
			&user.Status,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
