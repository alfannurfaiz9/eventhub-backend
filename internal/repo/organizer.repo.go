package repo

import (
	"context"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizerRepo struct {
	db *pgxpool.Pool
}

func NewOrganizerRepo(db *pgxpool.Pool) *OrganizerRepo {
	return &OrganizerRepo{
		db: db,
	}
}

func (o *OrganizerRepo) GetOrganizerDashboard(ctx context.Context, id int) (model.OrganizerDashboard, error) {
	sql := `
	SELECT COUNT(e.id), COUNT(ue.event_id), CONCAT(COUNT(e.id) * 100 / SUM(e.capacity), '%')
	FROM events e
	LEFT JOIN users u ON u.id = e.organizer_id
	LEFT JOIN user_event ue ON ue.event_id = e.id
	WHERE u.id = $1`
	args := []any{id}

	var data model.OrganizerDashboard
	err := o.db.QueryRow(ctx, sql, args...).Scan(
		&data.TotalEvent,
		&data.TotalAttendee,
		&data.AvgFillRate)

	if err != nil {
		return model.OrganizerDashboard{}, err
	}

	return data, nil
}

func (o *OrganizerRepo) GetOrganizerEvent(ctx context.Context, id int) ([]model.EventList, error) {
	sql := `
	SELECT e.title, e.img_url, e.start_at, l.name, e.capacity, COUNT(ue.user_id) 
	FROM events e
	LEFT JOIN users u ON u.id = e.organizer_id
	LEFT JOIN locations l ON l.id = e.location_id
	LEFT JOIN user_event ue ON ue.event_id = e.id
	WHERE u.id = $1
	GROUP BY e.id, u.id, l.id`
	args := []any{id}

	rows, err := o.db.Query(ctx, sql, args...)

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
