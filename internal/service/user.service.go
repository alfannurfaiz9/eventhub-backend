package service

import (
	"context"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
)

type UserService struct {
	ur *repo.UserRepo
}

func NewUserService(ur *repo.UserRepo) *UserService {
	return &UserService{
		ur: ur,
	}
}

func (u *UserService) GetUserProfile(ctx context.Context, id int) (dto.User, error) {
	result, err := u.ur.GetUserProfile(ctx, id)

	data := dto.User{
		FullName: result.FullName,
		Email:    result.Email,
		ImgUrl:   result.ImgUrl,
		Address:  result.Address,
		Bio:      result.Bio,
		Role:     result.Role,
	}

	return data, err
}

func (u *UserService) GetMyEvent(ctx context.Context, id int) ([]dto.EventList, error) {
	result, err := u.ur.GetMyEvent(ctx, id)

	data := make([]dto.EventList, 0, len(result))

	for _, v := range result {
		data = append(data, dto.EventList{
			Title:         v.Event.Title,
			ImgUrl:        v.Event.ImgUrl,
			StartAt:       v.Event.StartAt,
			Location:      v.Location.Name,
			Capacity:      v.Event.Capacity,
			TotalAttendee: v.TotalAttendee,
		})
	}

	return data, err
}

func (u *UserService) GetNotification(ctx context.Context, id int) ([]dto.Notification, error) {
	result, err := u.ur.GetNotification(ctx, id)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, custom_error.NotificationNotFound
	}

	data := make([]dto.Notification, 0, len(result))

	for _, v := range result {
		data = append(data, dto.Notification{
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt,
		})
	}

	return data, nil
}
