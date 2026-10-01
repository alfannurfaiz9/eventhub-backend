package controller

import (
	"errors"
	"net/http"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	us *service.UserService
}

func NewUserController(us *service.UserService) *UserController {
	return &UserController{
		us: us,
	}
}

func (u *UserController) GetUserProfile(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	result, err := u.us.GetUserProfile(ctx.Request.Context(), claims.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "success",
	})
}

func (u *UserController) GetMyEvent(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	result, err := u.us.GetMyEvent(ctx.Request.Context(), userId.(int))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "success",
	})
}

func (u *UserController) GetNotification(ctx *gin.Context) {
	id, _ := ctx.Get("id")
	result, err := u.us.GetNotification(ctx.Request.Context(), id.(int))

	if err != nil {
		if errors.Is(err, custom_error.NotificationNotFound) {
			ctx.JSON(http.StatusForbidden, dto.Response{
				Success: false,
				Message: err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get notification",
	})
}
