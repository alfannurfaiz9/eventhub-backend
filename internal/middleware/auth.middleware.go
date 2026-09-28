package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckToken(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "please login first",
		})
		return
	}

	result := strings.Split(bearer, " ")
	if len(result) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "invalid bearer token",
		})
		return
	}
	if result[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "invalid bearer token",
		})
		return
	}

	var token pkg.JWTClaims
	err := token.DecodeToken(result[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "invalid token",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "terjadi kesalahan sistem",
		})
		return
	}

	c.Set("token", token)
	c.Set("id", token.Id)
	c.Next()
}
