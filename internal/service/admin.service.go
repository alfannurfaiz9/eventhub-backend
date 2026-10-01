package service

import (
	"context"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
)

type AdminService struct {
	ar *repo.AdminRepo
}

func NewAdminService(ar *repo.AdminRepo) *AdminService {
	return &AdminService{
		ar: ar,
	}
}

func (a *AdminService) GetAdminDashboard(ctx context.Context) (dto.AdminDashboard, error) {
	result, err := a.ar.GetAdminDashboard(ctx)

	data := dto.AdminDashboard{
		TotalUser:      result.TotalUser,
		TotalEvent:     result.TotalEvent,
		TotalCommunity: result.TotalCommunity,
		AvgFillRate:    result.AvgFillRate,
	}

	return data, err
}

func (a *AdminService) GetAllUser(ctx context.Context) ([]dto.UserList, error) {
	result, err := a.ar.GetAllUser(ctx)

	userlist := make([]dto.UserList, 0, len(result))

	for _, v := range result {
		userlist = append(userlist, dto.UserList{
			FullName:  v.FullName,
			Email:     v.Email,
			Role:      v.Role,
			Status:    v.Status,
			CreatedAt: v.CreatedAt,
		})
	}

	return userlist, err
}
