package service

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/jackc/pgx/v5"
)

type CommunityService struct {
	cr *repo.CommunityRepo
}

func NewCommunityService(cr *repo.CommunityRepo) *CommunityService {
	return &CommunityService{
		cr: cr,
	}
}

func (c *CommunityService) GetCommunities(ctx context.Context, categories string, page int) ([]dto.CommunityList, error) {
	if page < 1 {
		return nil, custom_error.EventErrorPage
	}

	result, err := c.cr.GetCommunities(ctx, categories, page)

	if err != nil {
		return nil, err
	}

	data := make([]dto.CommunityList, 0, len(result))

	for _, v := range result {
		data = append(data, dto.CommunityList{
			Name:          v.Community.Name,
			Image:         v.Community.ImgUrl,
			Description:   v.Community.Description,
			Category:      v.Category.Name,
			TotalMember:   v.TotalMember,
			UpcomingEvent: v.UpcomingEvent,
		})
	}

	return data, nil
}

func (c *CommunityService) GetCommunityDetail(ctx context.Context, id int) (dto.CommunityList, error) {
	result, err := c.cr.GetCommunityDetail(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.CommunityList{}, custom_error.CommunityNotFound
		}
		return dto.CommunityList{}, err
	}

	data := dto.CommunityList{
		Name:          result.Community.Name,
		Image:         result.Community.ImgUrl,
		Description:   result.Community.Description,
		Category:      result.Category.Name,
		TotalMember:   result.TotalMember,
		UpcomingEvent: result.UpcomingEvent,
	}

	return data, nil
}

func (c *CommunityService) GetCommunityEvent(ctx context.Context, id int) ([]dto.EventList, error) {
	_, err := c.cr.GetCommunityDetail(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, custom_error.CommunityNotFound
		}
		return nil, err
	}

	result, err := c.cr.GetCommunityEvent(ctx, id)

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

func (c *CommunityService) GetCommunityMember(ctx context.Context, id int) ([]dto.CommunityMember, error) {
	_, err := c.cr.GetCommunityDetail(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, custom_error.CommunityNotFound
		}
		return nil, err
	}

	result, err := c.cr.GetCommunityMember(ctx, id)

	data := make([]dto.CommunityMember, 0, len(result))
	for _, v := range result {
		data = append(data, dto.CommunityMember{
			FullName: v.FullName,
		})
	}

	return data, err
}

func (c *CommunityService) GetPopularCommunity(ctx context.Context) ([]dto.CommunityList, error) {
	result, err := c.cr.GetPopularCommunity(ctx)

	data := make([]dto.CommunityList, 0, len(result))
	for _, v := range result {
		data = append(data, dto.CommunityList{
			Name:          v.Community.Name,
			Image:         v.Community.ImgUrl,
			Description:   v.Community.Description,
			Category:      v.Category.Name,
			TotalMember:   v.TotalMember,
			UpcomingEvent: v.UpcomingEvent,
		})
	}

	return data, err
}

func (c *CommunityService) JoinCommunity(ctx context.Context, userId, communityId int) error {
	return c.cr.JoinCommunity(ctx, userId, communityId)
}

func (c *CommunityService) LeaveCommunity(ctx context.Context, userId, communityId int) error {
	return c.cr.LeaveCommunity(ctx, userId, communityId)
}
