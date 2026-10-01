package repo

import (
	"context"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{
		db: db,
	}
}

func (a *AdminRepo) GetAdminDashboard(ctx context.Context) (model.AdminDashboard, error) {
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

func (a *AdminRepo) GetAllUser(ctx context.Context) ([]model.User, error) {
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
