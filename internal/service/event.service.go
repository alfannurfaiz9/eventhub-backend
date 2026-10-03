package service

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type EventService struct {
	er  *repo.EventRepo
	rdb *redis.Client
}

func NewEventService(er *repo.EventRepo, rdb *redis.Client) *EventService {
	return &EventService{
		er:  er,
		rdb: rdb,
	}
}

func (e *EventService) GetEvents(ctx context.Context, search, location, category string, page int) ([]dto.EventList, error) {
	if page == 0 {
		page += 1
	}

	if page < 0 {
		return nil, custom_error.EventErrorPage
	}

	result, err := e.er.GetEvents(ctx, search, location, category, page)

	data := make([]dto.EventList, 0, len(result))

	for _, v := range result {
		data = append(data, dto.EventList{
			Title:         v.Event.Title,
			ImgUrl:        v.Event.ImgUrl,
			Category:      v.Category.Name,
			StartAt:       v.Event.StartAt,
			Location:      v.Location.Name,
			TotalAttendee: v.TotalAttendee,
			Capacity:      v.Capacity,
		})
	}

	return data, err
}

func (e *EventService) GetEventDetail(ctx context.Context, id int) (dto.EventDetail, error) {
	result, err := e.er.GetEventDetail(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.EventDetail{}, custom_error.EventNotFound
		}
		return dto.EventDetail{}, err
	}

	data := dto.EventDetail{
		Title:         result.Event.Title,
		ImgUrl:        *result.Event.ImgUrl,
		Description:   result.Event.Description,
		Category:      result.Category.Name,
		StartAt:       result.Event.StartAt,
		Location:      result.Location.Name,
		TotalAttendee: result.TotalAttendee,
		Capacity:      result.Event.Capacity,
		Organizer:     result.User.FullName,
		Community:     result.Community.Name,
	}

	return data, nil
}

func (e *EventService) JoinEvent(ctx context.Context, userId int, event_id int) error {
	if err := e.er.JoinEvent(ctx, userId, event_id); err != nil {
		return err
	}

	return nil
}

func (e *EventService) GetUpcomingEvent(ctx context.Context) ([]dto.EventList, error) {
	result, err := e.er.GetUpcomingEvent(ctx)

	data := make([]dto.EventList, 0, len(result))

	for _, v := range result {
		data = append(data, dto.EventList{
			Title:         v.Event.Title,
			ImgUrl:        v.Event.ImgUrl,
			Category:      v.Category.Name,
			StartAt:       v.Event.StartAt,
			Location:      v.Location.Name,
			TotalAttendee: v.TotalAttendee,
			Capacity:      v.Capacity,
		})
	}

	return data, err
}

func (e *EventService) LeaveEvent(ctx context.Context, user_id int, eventId int) error {
	if err := e.er.LeaveEvent(ctx, user_id, eventId); err != nil {
		return err
	}

	return nil
}
