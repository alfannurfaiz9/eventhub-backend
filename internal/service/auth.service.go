package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	ar  *repo.AuthRepo
	rdb *redis.Client
}

func NewAuthService(ar *repo.AuthRepo, rdb *redis.Client) *AuthService {
	return &AuthService{
		ar:  ar,
		rdb: rdb,
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

func (a *AuthService) Logout(ctx context.Context, token string) error {
	redisToken, err := a.rdb.Get(ctx, "eventhub:blacklist_token").Result()
	if err != nil {
		log.Println(err)
	} else {
		var redisJson []string
		if err := json.Unmarshal([]byte(redisToken), &redisJson); err != nil {
			return err
		}

		splittedToken := strings.Split(token, " ")[1]
		redisJson = append(redisJson, splittedToken)
		tokenJson, err := json.Marshal(redisJson)
		if err != nil {
			return err
		}

		if err := a.rdb.Set(ctx, "eventhub:blacklist_token", tokenJson, 0).Err(); err != nil {
			return err
		}

		return nil
	}

	splittedToken := strings.Split(token, " ")[1]
	tokenArr := []string{splittedToken}
	tokenJson, err := json.Marshal(tokenArr)
	if err != nil {
		return err
	}

	if err := a.rdb.Set(ctx, "eventhub:blacklist_token", tokenJson, 5*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ForgotPassword(ctx context.Context, body dto.ForgotPassword) error {
	if len(body.Email) == 0 {
		return custom_error.UserNotFound
	}

	if len(body.NewPassword) < 6 {
		return custom_error.ForgotPasswordInvalidLength
	}

	email, err := a.ar.FindUserByEmail(ctx, model.User{
		Email: body.Email,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return custom_error.UserNotFound
		}
		return err
	}

	hash := pkg.NewRecommendHashConfig()
	hashedPass := hash.GenerateHash(body.NewPassword)

	if err := a.ar.ForgotPassword(ctx, model.User{
		Email:    email,
		Password: hashedPass,
	}); err != nil {
		return err
	}

	return nil
}
