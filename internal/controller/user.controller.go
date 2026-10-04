package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
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

// GetUserProfile
//
// @Summary			Get user profile
// @Description		Get user profile detail
// @Tags			user
// @Produce			json
// @Router			/user/profile		[get]
// @Security		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (u *UserController) GetUserProfile(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	result, err := u.us.GetUserProfile(ctx.Request.Context(), claims.Id)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get user profile",
	})
}

func (u *UserController) GetMyEvent(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	result, err := u.us.GetMyEvent(ctx.Request.Context(), claims.Id)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get my event",
	})
}

func (u *UserController) GetNotification(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	result, err := u.us.GetNotification(ctx.Request.Context(), claims.Id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.NotificationNotFound) {
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
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
		Data:    result,
		Message: "successfully get notification",
	})
}

// ChangeUserProfile
//
// @Summary			Change user profile
// @Description		Change user profile
// @Tags			user
// @Accept			mpfd
// @Produce			json
// @Security		BearerToken
// @Router			/user/edit			[patch]
// @Param			full_name			formData	string	false	"full name"
// @Param			img_url				formData	file	false	"img url"
// @Param			address				formData	string	false	"address"
// @Param			bio					formData	string	false	"bio"
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			406		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (u *UserController) ChangeUserProfile(ctx *gin.Context) {
	var body dto.UpdateProfile

	if err := ctx.ShouldBindWith(&body, binding.FormMultipart); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	var imgUrl string

	if body.Img != nil {
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), body.FullName, path.Ext(body.Img.Filename))
		filepath := path.Join("public", "img", filename)

		if err := ctx.SaveUploadedFile(body.Img, filepath); err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Message: "internal server error",
			})

			return
		}

		savedPath := fmt.Sprintf("img/%s", filename)
		imgUrl = savedPath
	}

	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	if body.Img == nil {
		imgUrl = ""
	}

	if err := u.us.ChangeUserProfile(ctx.Request.Context(), body, claims.Id, &imgUrl); err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.NoRowsAffected) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})

			return
		}

		if errors.Is(err, custom_error.AllFieldIsEmpty) {
			ctx.JSON(http.StatusNotAcceptable, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})

			return
		}

		if errors.Is(err, custom_error.AllFieldIsEmpty) {
			ctx.JSON(http.StatusNotAcceptable, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})

			return
		}

		if errors.Is(err, custom_error.ChangeUserrInvalidLength) {
			ctx.JSON(http.StatusNotAcceptable, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})

			return
		}

		if errors.Is(err, custom_error.ChangeUserInvalidPassword) {
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

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "profile updated successfully",
	})
}

// ChangeUserPassword
//
// @Summary			Change user password
// @Description		Change user password
// @Tags			user
// @Accept			json
// @Security		BearerToken
// @Router			/user/change-password	[post]
// @Param			data	body		dto.ChangePassword	true	"body to change password"
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			406		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (u *UserController) ChangeUserPassword(ctx *gin.Context) {
	var body dto.ChangePassword

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	if err := u.us.ChangeUserPassword(ctx.Request.Context(), body, claims.Id); err != nil {
		if errors.Is(err, custom_error.ChangeUserrInvalidLength) {
			ctx.JSON(http.StatusNotAcceptable, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})

			return
		}

		if errors.Is(err, custom_error.ChangeUserInvalidPassword) {
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
		Message: "password successfully changed",
	})
}

// GetUserInformation
//
// @Summary			Get user information
// @Description		Get user information for header
// @Tags			user
// @Produce			json
// @Router			/user/info		[get]
// @Security		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (u *UserController) GetUserInformation(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)
	result, err := u.us.GetUserInformation(ctx.Request.Context(), claims.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get user information",
	})
}
