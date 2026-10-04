package service

import (
	"context"
	"errors"
	"log"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type EventService struct {
	er  *repo.EventRepo
	rdb *redis.Client
	db  *pgxpool.Pool
}

func NewEventService(er *repo.EventRepo, rdb *redis.Client, db *pgxpool.Pool) *EventService {
	return &EventService{
		er:  er,
		rdb: rdb,
		db:  db,
	}
}

func (e *EventService) GetEvents(ctx context.Context, search, location, category string, page int) ([]dto.EventList, error) {
	if page < 1 {
		return nil, custom_error.EventErrorPage
	}

	result, err := e.er.GetEvents(ctx, search, location, category, page, e.db)

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
	result, err := e.er.GetEventDetail(ctx, id, e.db)

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
	if err := e.er.JoinEvent(ctx, userId, event_id, e.db); err != nil {
		return err
	}

	return nil
}

func (e *EventService) GetUpcomingEvent(ctx context.Context) ([]dto.EventList, error) {
	result, err := e.er.GetUpcomingEvent(ctx, e.db)

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
	if err := e.er.LeaveEvent(ctx, user_id, eventId, e.db); err != nil {
		return err
	}

	return nil
}

func (e *EventService) SaveEvent(ctx context.Context, userId, eventId int) error {
	return e.er.SaveEvent(ctx, userId, eventId, e.db)
}

func (e *EventService) CreateEvent(ctx context.Context, body dto.CreateEvent, organizerId int) error {
	tx, err := e.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err)
		}
	}()

	speakerId, err := e.er.InsertSpeakers(ctx, model.Speaker{
		Name:     body.SpeakerName,
		ImgUrl:   body.SpeakerImgUrl,
		Position: body.SpeakerPosition,
		Company:  body.SpeakerCompany,
	}, tx)
	if err != nil {
		return err
	}

	locationId, err := e.er.InsertLocation(ctx, model.Location{Name: body.LocationName}, tx)
	if err != nil {
		return err
	}

	eventId, err := e.er.InsertEvent(ctx, model.Event{
		Title:       body.Title,
		ImgUrl:      body.ImgUrl,
		Description: body.Description,
		StartAt:     body.StartAt,
		EndAt:       body.EndAt,
		Format:      body.Format,
		Capacity:    body.Capacity,
		OrganizerId: organizerId,
		CommunityId: body.CommunityId,
		LocationId:  locationId,
	}, tx)
	if err != nil {
		return err
	}

	cmdEc, err := e.er.InsertEventCategory(ctx, model.EventCategory{
		EventId:    eventId,
		CategoryId: body.CategoryId,
	}, tx)

	if err != nil {
		return err
	}

	if cmdEc.RowsAffected() == 0 {
		return custom_error.NoRowsAffected
	}

	cmdEs, err := e.er.InsertEventSpeaker(ctx, model.EventSpeaker{
		EventId:   eventId,
		SpeakerId: speakerId,
	}, tx)

	if err != nil {
		return err
	}

	if cmdEs.RowsAffected() == 0 {
		return custom_error.NoRowsAffected
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println(err)
	}

	return nil
}
