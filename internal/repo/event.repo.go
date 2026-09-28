package repo

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepo struct {
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) *EventRepo {
	return &EventRepo{
		db: db,
	}
}

func (e *EventRepo) GetEvents(ctx context.Context, search, location, category string) ([]model.EventList, error) {
	sql := `
	SELECT events.title, events.img_url, STRING_AGG(categories.name, ', '), events.start_at, locations.name, COUNT(user_event.event_id), events.capacity 
	FROM events 
	LEFT JOIN locations ON locations.id = events.location_id 
	LEFT JOIN communities ON communities.id = events.community_id 
	LEFT JOIN event_category ON event_category.event_id = events.id 
	LEFT JOIN categories ON categories.id = event_category.category_id 
	LEFT JOIN user_event ON user_event.event_id = events.id 
	WHERE events.title ILIKE $1 AND locations.name ILIKE $2
	GROUP BY events.id, categories.id, locations.id
	HAVING STRING_AGG(categories.name, ', ') ILIKE $3`
	args := []any{"%" + search + "%", "%" + location + "%", "%" + category + "%"}

	rows, err := e.db.Query(ctx, sql, args...)

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
			&event.Event.Capacity,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}

func (e *EventRepo) GetEventDetail(ctx context.Context, id string) (model.EventDetail, error) {
	sql := "SELECT events.title, events.img_url, events.description, categories.name, events.start_at, locations.name, COUNT(user_event.event_id), events.capacity, users.full_name, communities.name FROM events LEFT JOIN locations ON locations.id = events.location_id LEFT JOIN communities ON communities.id = events.community_id LEFT JOIN event_category ON event_category.event_id = events.id LEFT JOIN categories ON categories.id = event_category.category_id LEFT JOIN user_event ON user_event.event_id = events.id LEFT JOIN users ON users.id = events.organizer_id WHERE events.id = $1 GROUP BY events.id, categories.id, locations.id, users.id, communities.id"
	args := []any{id}

	var data model.EventDetail
	if err := e.db.QueryRow(ctx, sql, args...).Scan(
		&data.Event.Title,
		&data.Event.ImgUrl,
		&data.Event.Description,
		&data.Category.Name,
		&data.Event.StartAt,
		&data.Location.Name,
		&data.TotalAttendee,
		&data.Event.Capacity,
		&data.User.FullName,
		&data.Community.Name,
	); err != nil {
		return model.EventDetail{}, err
	}

	return data, nil
}

func (e *EventRepo) JoinEvent(ctx context.Context, user_id int, body model.UserEvent) error {
	sql := "INSERT INTO user_event(user_id, event_id) VALUES($1, $2)"
	args := []any{user_id, body.EventId}

	cmd, err := e.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}

func (e *EventRepo) GetUpcomingEvent(ctx context.Context) ([]model.EventList, error) {
	sql := `
	SELECT events.title, events.img_url, STRING_AGG(categories.name, ', '), events.start_at, locations.name, COUNT(user_event.event_id), events.capacity 
	FROM events 
	LEFT JOIN locations ON locations.id = events.location_id 
	LEFT JOIN communities ON communities.id = events.community_id 
	LEFT JOIN event_category ON event_category.event_id = events.id 
	LEFT JOIN categories ON categories.id = event_category.category_id 
	LEFT JOIN user_event ON user_event.event_id = events.id 
	WHERE events.start_at > NOW()
	GROUP BY events.id, categories.id, locations.id`

	rows, err := e.db.Query(ctx, sql)

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
			&event.Event.Capacity,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}

func (e *EventRepo) LeaveEvent(ctx context.Context, userId int, body model.UserEvent) error {
	sql := "DELETE FROM user_event WHERE user_id = $1 AND event_id = $2"
	args := []any{userId, body.EventId}

	cmd, err := e.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no row affected")

	}

	return nil
}
