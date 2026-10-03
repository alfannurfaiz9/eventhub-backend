package service

import (
	"context"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
)

type UserService struct {
	ur *repo.UserRepo
}

func NewUserService(ur *repo.UserRepo) *UserService {
	return &UserService{
		ur: ur,
	}
}

func (u *UserService) GetUserProfile(ctx context.Context, id int) (dto.UserProfile, error) {
	result, err := u.ur.GetUserProfile(ctx, id)

	data := dto.UserProfile{
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

func (u *UserService) ChangeUserProfile(ctx context.Context, body dto.UpdateProfile, id int) error {
	if len(body.NewPassword) < 6 && len(body.NewPassword) != 0 {
		return custom_error.ChangeUserrInvalidLength
	}

	currentProfile, err := u.ur.FindUser(ctx, id)

	if err != nil {
		return err
	}

	if err := pkg.CompareHash(body.CurrentPassword, currentProfile.Password); err != nil {
		return custom_error.ChangeUserInvalidPassword
	}

	hash := pkg.NewRecommendHashConfig()
	hashedPass := hash.GenerateHash(body.NewPassword)

	fixNewPassword := hashedPass
	if len(body.NewPassword) == 0 {
		fixNewPassword = currentProfile.Password
	}

	if err := u.ur.ChangeUserProfile(ctx, model.User{
		FullName: body.FullName,
		ImgUrl:   body.ImgUrl,
		Address:  body.Address,
		Bio:      body.Bio,
		Password: fixNewPassword},
		id); err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUserInformation(ctx context.Context, id int) (dto.UserInformation, error) {
	result, err := u.ur.GetUserInformation(ctx, id)

	data := dto.UserInformation{
		FullName: result.FullName,
		Email:    result.Email,
		ImgUrl:   result.ImgUrl,
	}

	return data, err
}
