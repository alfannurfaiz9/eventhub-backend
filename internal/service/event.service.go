package service

import (
	"context"
	"errors"
	"fmt"
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

func (e *EventService) CreateEvent(ctx context.Context, body dto.CreateEvent, organizerId int, eventImg string) error {
	tx, err := e.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err)
		}
	}()

	speakersId := make([]int, 0, len(body.SpeakerName))
	for i := range body.SpeakerName {
		speakerId, err := e.er.InsertSpeakers(ctx, model.Speaker{
			Name:     body.SpeakerName[i],
			Position: body.SpeakerPosition[i],
			Company:  body.SpeakerCompany[i],
		}, tx)
		if err != nil {
			return err
		}

		speakersId = append(speakersId, speakerId)
	}

	locationId, err := e.er.InsertLocation(ctx, model.Location{Name: body.LocationName}, tx)
	if err != nil {
		return err
	}

	eventId, err := e.er.InsertEvent(ctx, model.Event{
		Title:       body.Title,
		ImgUrl:      &eventImg,
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

	catIds := make([]model.EventCategory, 0, len(body.CategoryId))
	for _, v := range body.CategoryId {
		catIds = append(catIds, model.EventCategory{
			EventId:    eventId,
			CategoryId: v,
		})
	}

	cmdEc, err := e.er.InsertEventCategory(ctx, catIds, tx)

	if err != nil {
		fmt.Println(catIds)

		return err
	}

	if cmdEc.RowsAffected() == 0 {
		return custom_error.NoRowsAffected
	}

	speakers := make([]model.EventSpeaker, 0, len(speakersId))
	for _, v := range speakersId {
		speakers = append(speakers, model.EventSpeaker{
			EventId:   eventId,
			SpeakerId: v,
		})
	}
	cmdEs, err := e.er.InsertEventSpeaker(ctx, speakers, tx)

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
