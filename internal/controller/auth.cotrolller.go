package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	as *service.AuthService
}

func NewAuthController(as *service.AuthService) *AuthController {
	return &AuthController{
		as: as,
	}
}

func (a *AuthController) Register(ctx *gin.Context) {
	var body dto.User

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	if err := a.as.Register(ctx.Request.Context(), body); err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "account created successfully",
	})
}

func (a *AuthController) Login(ctx *gin.Context) {
	var body dto.User

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	token, err := a.as.Login(ctx.Request.Context(), body)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.EmptyLoginField) {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		if errors.Is(err, custom_error.LoginInvalidEmailOrPassword) {
			ctx.JSON(http.StatusUnauthorized, dto.Response{
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

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data: gin.H{
			"token": token,
		},
		Message: "login success",
	})
}

func (a *AuthController) GetUserProfile(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	result, err := a.as.GetUserProfile(ctx.Request.Context(), claims.Id)

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

func (a *AuthController) ChangeUserPassword(ctx *gin.Context) {
	var body dto.User

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	if err := a.as.ChangeUserPassword(ctx, body, claims.Id); err != nil {
		ctx.JSON(http.StatusNotAcceptable, dto.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "password successfully changed",
	})
}

func (a *AuthController) GetMyEvent(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	result, err := a.as.GetMyEvent(ctx.Request.Context(), userId.(int))

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
