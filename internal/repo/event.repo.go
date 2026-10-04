package repo

import (
	"context"
	"errors"

	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type EventRepo struct{}

func NewEventRepo() *EventRepo {
	return &EventRepo{}
}

func (e *EventRepo) GetEvents(ctx context.Context, search, location, category string, page int, db DBTX) ([]model.EventList, error) {
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
	HAVING STRING_AGG(categories.name, ', ') ILIKE $3
	LIMIT $4 OFFSET $5`
	limit := 6
	offset := limit * (page - 1)
	args := []any{"%" + search + "%", "%" + location + "%", "%" + category + "%", limit, offset}

	rows, err := db.Query(ctx, sql, args...)

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

func (e *EventRepo) GetEventDetail(ctx context.Context, id int, db DBTX) (model.EventDetail, error) {
	sql := "SELECT events.title, events.img_url, events.description, categories.name, events.start_at, locations.name, COUNT(user_event.event_id), events.capacity, users.full_name, communities.name FROM events LEFT JOIN locations ON locations.id = events.location_id LEFT JOIN communities ON communities.id = events.community_id LEFT JOIN event_category ON event_category.event_id = events.id LEFT JOIN categories ON categories.id = event_category.category_id LEFT JOIN user_event ON user_event.event_id = events.id LEFT JOIN users ON users.id = events.organizer_id WHERE events.id = $1 GROUP BY events.id, categories.id, locations.id, users.id, communities.id"
	args := []any{id}

	var data model.EventDetail
	if err := db.QueryRow(ctx, sql, args...).Scan(
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

func (e *EventRepo) JoinEvent(ctx context.Context, user_id int, event_id int, db DBTX) error {
	sql := "INSERT INTO user_event(user_id, event_id) VALUES($1, $2)"
	args := []any{user_id, event_id}

	cmd, err := db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}

func (e *EventRepo) GetUpcomingEvent(ctx context.Context, db DBTX) ([]model.EventList, error) {
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

	rows, err := db.Query(ctx, sql)

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

func (e *EventRepo) LeaveEvent(ctx context.Context, userId int, eventId int, db DBTX) error {
	sql := "DELETE FROM user_event WHERE user_id = $1 AND event_id = $2"
	args := []any{userId, eventId}

	cmd, err := db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no row affected")

	}

	return nil
}

func (e *EventRepo) SaveEvent(ctx context.Context, userId, eventId int, db DBTX) error {
	sql := `
	INSERT INTO user_saved_event
	VALUES($1, $2)`
	args := []any{userId, eventId}

	cmd, err := db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return custom_error.NoRowsAffected
	}

	return nil
}

func (e *EventRepo) InsertSpeakers(ctx context.Context, body model.Speaker, db DBTX) (int, error) {
	sql := `
	INSERT INTO speakers(name, img_url, position, company)
	VALUES ($1, $2, $3, $4) 
	RETURNING id`
	args := []any{body.Name, body.ImgUrl, body.Position, body.Company}

	var id int
	if err := db.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (e *EventRepo) InsertLocation(ctx context.Context, body model.Location, db DBTX) (int, error) {
	sql := `
	INSERT INTO locations(name)
	VALUES($1) 
	RETURNING id`
	args := []any{body.Name}

	var id int
	if err := db.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (e *EventRepo) InsertEvent(ctx context.Context, body model.Event, db DBTX) (int, error) {
	sql := `
	INSERT INTO events(title, img_url, description, start_at, end_at, format, capacity, organizer_id, community_id, location_id)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	RETURNING id`
	args := []any{body.Title, body.ImgUrl, body.Description, body.StartAt, body.EndAt, body.Format, body.Capacity, body.OrganizerId, body.CommunityId, body.LocationId}

	var id int
	if err := db.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (e *EventRepo) InsertEventCategory(ctx context.Context, body model.EventCategory, db DBTX) (pgconn.CommandTag, error) {
	sql := `
	INSERT INTO event_category(event_id, category_id)
	VALUES($1, $2)
	`
	args := []any{body.EventId, body.CategoryId}

	return db.Exec(ctx, sql, args...)
}

func (e *EventRepo) InsertEventSpeaker(ctx context.Context, body model.EventSpeaker, db DBTX) (pgconn.CommandTag, error) {
	sql := `
	INSERT INTO event_speakers(event_id, speaker_id)
	VALUES($1, $2)
	`
	args := []any{body.EventId, body.SpeakerId}

	return db.Exec(ctx, sql, args...)
}
