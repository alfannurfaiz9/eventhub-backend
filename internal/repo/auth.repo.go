package repo

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
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

func (a *AuthRepo) GetUserProfile(ctx context.Context, id int) (model.User, error) {
	sql := "SELECT full_name, email, img_url, address, bio, role from users WHERE id = $1"
	args := []any{id}

	var user model.User
	if err := a.db.QueryRow(ctx, sql, args...).Scan(&user.FullName, &user.Email, &user.ImgUrl, &user.Address, &user.Bio, &user.Role); err != nil {
		return model.User{}, err
	}

	return user, nil
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

func (a *AuthRepo) GetMyEvent(ctx context.Context, id int) ([]model.EventList, error) {
	sql := `
	SELECT events.title, events.img_url, STRING_AGG(categories.name, ', '), events.start_at, locations.name, COUNT(user_event.event_id), events.capacity 
	FROM events 
	LEFT JOIN locations ON locations.id = events.location_id 
	LEFT JOIN communities ON communities.id = events.community_id 
	LEFT JOIN event_category ON event_category.event_id = events.id 
	LEFT JOIN categories ON categories.id = event_category.category_id 
	LEFT JOIN user_event ON user_event.event_id = events.id 
	WHERE user_event.user_id = $1
	GROUP BY events.id, categories.id, locations.id`
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
			&event.Category.Name,
			&event.Event.StartAt,
			&event.Location.Name,
			&event.TotalAttendee,
			&event.Event.Capacity); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (a *AuthRepo) GetNotification(ctx context.Context, id int) ([]model.Notification, error) {
	sql := `
	SELECT notifications.title, notifications.description, notifications.created_at
	FROM notifications
	LEFT JOIN users ON users.id = notifications.user_id
	WHERE users.id = $1`
	args := []any{id}

	rows, err := a.db.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	var notifications []model.Notification

	for rows.Next() {
		var notification model.Notification

		if err := rows.Scan(
			&notification.Title,
			&notification.Description,
			&notification.CreatedAt,
		); err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (a *AuthRepo) GetOrganizerDashboard(ctx context.Context, id int) (dto.OrganizerDashboard, error) {
	sql := `
	SELECT COUNT(e.id), COUNT(ue.event_id), CONCAT(COUNT(e.id) * 100 / SUM(e.capacity), '%')
	FROM events e
	LEFT JOIN users u ON u.id = e.organizer_id
	LEFT JOIN user_event ue ON ue.event_id = e.id
	WHERE u.id = $1`
	args := []any{id}

	var data dto.OrganizerDashboard
	err := a.db.QueryRow(ctx, sql, args...).Scan(
		&data.TotalEvent,
		&data.TotalAttendee,
		&data.AvgFillRate)

	if err != nil {
		return dto.OrganizerDashboard{}, err
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
