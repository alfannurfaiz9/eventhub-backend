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

func (a *AuthRepo) ChangeUserPassword(ctx context.Context, body model.User, id int) error {
	sql := "UPDATE users SET password = $1 WHERE id = $2"
	args := []any{body.Password, id}

	cmd, err := a.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}

func (a *AuthRepo) GetOrganizerDashboard(ctx context.Context, id int) (model.OrganizerDashboard, error) {
	sql := `
	SELECT COUNT(e.id), COUNT(ue.event_id), CONCAT(COUNT(e.id) * 100 / SUM(e.capacity), '%')
	FROM events e
	LEFT JOIN users u ON u.id = e.organizer_id
	LEFT JOIN user_event ue ON ue.event_id = e.id
	WHERE u.id = $1`
	args := []any{id}

	var data model.OrganizerDashboard
	err := a.db.QueryRow(ctx, sql, args...).Scan(
		&data.TotalEvent,
		&data.TotalAttendee,
		&data.AvgFillRate)

	if err != nil {
		return model.OrganizerDashboard{}, err
	}

	return data, nil
}

func (a *AuthRepo) GetOrganizerEvent(ctx context.Context, id int) ([]model.EventList, error) {
	sql := `
	SELECT e.title, e.img_url, e.start_at, l.name, e.capacity, COUNT(ue.user_id) 
	FROM events e
	LEFT JOIN users u ON u.id = e.organizer_id
	LEFT JOIN locations l ON l.id = e.location_id
	LEFT JOIN user_event ue ON ue.event_id = e.id
	WHERE u.id = $1
	GROUP BY e.id, u.id, l.id`
	args := []any{id}

	rows, err := a.db.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	var events []model.EventList
	for rows.Next() {
		var event model.EventList
		if err := rows.Scan(
			&event.Event.Title,
			&event.Event.ImgUrl,
			&event.Event.StartAt,
			&event.Location.Name,
			&event.Event.Capacity,
			&event.TotalAttendee,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, err
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
