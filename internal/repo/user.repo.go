package repo

import (
	"context"

	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) GetUserProfile(ctx context.Context, id int) (model.User, error) {
	sql := "SELECT full_name, email, img_url, address, bio, role from users WHERE id = $1"
	args := []any{id}

	var user model.User
	if err := u.db.QueryRow(ctx, sql, args...).Scan(&user.FullName, &user.Email, &user.ImgUrl, &user.Address, &user.Bio, &user.Role); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *UserRepo) GetMyEvent(ctx context.Context, id int) ([]model.EventList, error) {
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

	rows, err := u.db.Query(ctx, sql, args...)

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

func (u *UserRepo) GetNotification(ctx context.Context, id int) ([]model.Notification, error) {
	sql := `
	SELECT notifications.title, notifications.description, notifications.created_at
	FROM notifications
	LEFT JOIN users ON users.id = notifications.user_id
	WHERE users.id = $1`
	args := []any{id}

	rows, err := u.db.Query(ctx, sql, args...)

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

func (u *UserRepo) FindUser(ctx context.Context, id int) (model.User, error) {
	sql := "SELECT id, email, password, role FROM users WHERE id=$1"
	args := []any{id}

	var data model.User
	if err := u.db.QueryRow(ctx, sql, args...).Scan(&data.Id, &data.Email, &data.Password, &data.Role); err != nil {
		return model.User{}, err
	}

	return data, nil
}

func (u *UserRepo) ChangeUserProfile(ctx context.Context, body model.User, id int) error {
	sql := `
	UPDATE users
	SET full_name = COALESCE(NULLIF($1, ''), full_name), img_url = COALESCE(NULLIF($2, ''), img_url), address = COALESCE(NULLIF($3, ''), address), bio = COALESCE(NULLIF($4, ''), bio), password = COALESCE(NULLIF($5, ''), password)
	WHERE id = $6`
	args := []any{body.FullName, body.ImgUrl, body.Address, body.Bio, body.Password, id}

	cmd, err := u.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return custom_error.NoRowsAffected
	}

	return nil
}

func (u *UserRepo) GetUserInformation(ctx context.Context, id int) (model.User, error) {
	sql := `
	SELECT full_name, email, img_url
	FROM users
	where id = $1`
	args := []any{id}

	var user model.User
	if err := u.db.QueryRow(ctx, sql, args...).Scan(
		&user.FullName,
		&user.Email,
		&user.ImgUrl,
	); err != nil {
		return model.User{}, err
	}

	return user, nil
}
