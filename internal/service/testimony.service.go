package service

import (
	"context"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
)

type TestimonyService struct {
	tr *repo.TestimonyRepo
}

func NewTestimonyService(tr *repo.TestimonyRepo) *TestimonyService {
	return &TestimonyService{
		tr: tr,
	}
}

func (t *TestimonyService) GetTestimony(ctx context.Context) ([]dto.Testimony, error) {
	result, err := t.tr.GetTestimony(ctx)

	data := make([]dto.Testimony, 0, len(result))

	for _, v := range result {
		data = append(data, dto.Testimony{
			Id:       v.Testimony.Id,
			Name:     v.User.FullName,
			Company:  v.Testimony.Company,
			Position: v.Testimony.Position,
			Message:  v.Testimony.Message,
		})
	}

	return data, err
}

func (t *TestimonyService) SetTestimony(ctx context.Context, userId int, body dto.Testimony) error {
	if body.Message == "" || body.Company == "" || body.Position == "" {
		return custom_error.TestimonyEmptyField
	}

	if err := t.tr.SetTestimony(ctx, userId, model.Testimony{Message: body.Message, Company: body.Company, Position: body.Position}); err != nil {
		return err
	}

	return nil
}
