package middleware

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleWare struct {
	rdb *redis.Client
}

func NewAuthMiddleWare(rdb *redis.Client) *AuthMiddleWare {
	return &AuthMiddleWare{
		rdb: rdb,
	}
}

func (a *AuthMiddleWare) UserMiddleware(ctx *gin.Context) {
	bearer := ctx.GetHeader("Authorization")
	if bearer == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "please login first",
		})
		return
	}

	result := strings.Split(bearer, " ")
	if len(result) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "invalid bearer token",
		})
		return
	}
	if result[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "invalid bearer token",
		})
		return
	}

	blacklistRedis, err := a.rdb.Get(ctx, "alfan:token").Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Println("redis key not exist")
		}
		log.Println(err.Error())
	}

	var blacklist []string
	json.Unmarshal([]byte(blacklistRedis), &blacklist)

	if slices.Contains(blacklist, result[1]) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "please login first",
		})

		return
	}

	var token pkg.JWTClaims
	err = token.DecodeToken(result[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: "invalid token",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})
		return
	}

	ctx.Set("token", token)
	ctx.Set("id", token.Id)
	ctx.Next()
}
