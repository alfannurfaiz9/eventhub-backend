package service

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	ar *repo.AuthRepo
}

func NewAuthService(ar *repo.AuthRepo) *AuthService {
	return &AuthService{
		ar: ar,
	}
}

func (a *AuthService) Register(ctx context.Context, body dto.Register) error {
	if len(body.Email) < 6 || len(body.Password) < 6 {
		return custom_error.RegisterInvalidLength
	}

	_, err := a.ar.FindUser(ctx, body.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if err == nil {
		return custom_error.RegisterAlreadyExist
	}

	hash := pkg.NewRecommendHashConfig()
	hashedPass := hash.GenerateHash(body.Password)

	if err := a.ar.Register(ctx, model.User{
		FullName: body.FullName,
		Email:    body.Email,
		Password: hashedPass,
	}); err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Login(ctx context.Context, body dto.Login) (string, error) {
	if len(body.Email) == 0 || len(body.Password) == 0 {
		return "", custom_error.EmptyLoginField
	}

	user, err := a.ar.FindUser(ctx, body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", custom_error.LoginInvalidEmailOrPassword
		}
		return "", err
	}

	if err := pkg.CompareHash(body.Password, user.Password); err != nil {
		return "", custom_error.LoginInvalidEmailOrPassword
	}

	claims := pkg.NewJWTClaims(user.Id, user.Role)
	return claims.GenToken()
}
