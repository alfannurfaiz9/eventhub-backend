package service

import (
	"context"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
)

type OrganizerService struct {
	or *repo.OrganizerRepo
}

func NewOrganizerService(or *repo.OrganizerRepo) *OrganizerService {
	return &OrganizerService{
		or: or,
	}
}

func (o *OrganizerService) GetOrganizerDashboard(ctx context.Context, id int) (dto.OrganizerDashboard, error) {
	result, err := o.or.GetOrganizerDashboard(ctx, id)

	data := dto.OrganizerDashboard{
		TotalEvent:    result.TotalAttendee,
		TotalAttendee: result.TotalAttendee,
		AvgFillRate:   result.AvgFillRate,
	}

	return data, err
}

func (o *OrganizerService) GetOrganizerEvent(ctx context.Context, id int) ([]dto.EventList, error) {
	result, err := o.or.GetOrganizerEvent(ctx, id)

	events := make([]dto.EventList, 0, len(result))
	for _, v := range result {
		events = append(events, dto.EventList{
			Title:         v.Event.Title,
			ImgUrl:        v.Event.ImgUrl,
			StartAt:       v.Event.StartAt,
			Location:      v.Location.Name,
			TotalAttendee: v.TotalAttendee,
		})
	}

	return events, err
}
