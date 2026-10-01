package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	id, _ := ctx.Get("id")

	result, err := u.us.GetUserProfile(ctx.Request.Context(), id.(int))

	if err != nil {
		log.Println(err.Error())
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
	id, _ := ctx.Get("id")

	result, err := u.us.GetMyEvent(ctx.Request.Context(), id.(int))

	if err != nil {
		log.Println(err.Error())
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
		log.Println(err.Error())
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

func (u *UserController) ChangeUserPassword(ctx *gin.Context) {
	var body dto.User

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	id, _ := ctx.Get("id")

	if err := u.us.ChangeUserPassword(ctx.Request.Context(), body, id.(int)); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusNotAcceptable, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "password successfully changed",
	})
}
