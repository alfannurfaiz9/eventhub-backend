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

type AuthController struct {
	as *service.AuthService
}

func NewAuthController(as *service.AuthService) *AuthController {
	return &AuthController{
		as: as,
	}
}

// Register
//
// @Summary			Create new user
// @Description		Create new user
// @Tags			auth
// @Produce			json
// @Router			/auth/register	[post]
// @Param			data	body	dto.Register	true	"body to register"
// @Success			201		{object}	dto.Response
// @Failure			400		{object}	dto.ErrorResponse
// @Failure			406		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (a *AuthController) Register(ctx *gin.Context) {
	var body dto.Register

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	if err := a.as.Register(ctx.Request.Context(), body); err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.RegisterInvalidLength) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})

			return
		}
		if errors.Is(err, custom_error.RegisterAlreadyExist) {
			ctx.JSON(http.StatusNotAcceptable, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})

			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "account created successfully",
	})
}

// Login
//
// @Summary			Login user
// @Description		Login user
// @Tags			auth
// @Accept			json
// @Produce			json
// @Router			/auth/login	[post]
// @Param			data	body	dto.Login	true	"body to login"
// @Success			200		{object}	dto.Response
// @Failure			400		{object}	dto.ErrorResponse
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (a *AuthController) Login(ctx *gin.Context) {
	var body dto.Login

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	token, err := a.as.Login(ctx.Request.Context(), body)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.EmptyLoginField) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		if errors.Is(err, custom_error.LoginInvalidEmailOrPassword) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"token": token,
		},
		Message: "login success",
	})
}
