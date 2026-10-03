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
// @Router			/user/change-profile		[patch]
// @Security		BearerToken
// @Param			full_name			formData	string	true	"full_name to change profile"
// @Param			img_url				formData	file	true	"img_url to change profile"
// @Param			address				formData	String	true	"address to change profile"
// @Param			bio					formData	String	true	"bio to change profile"
// @Param			current_password	formData	String	true	"current_password to change profile"
// @Param			new_password		formData	String	true	"new_password to change profile"
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

	log.Println(body.ImgUrl.Size)
	log.Println(body.ImgUrl.Filename)

	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), body.FullName, path.Ext(body.ImgUrl.Filename))
	filepath := path.Join("public", "img", filename)

	if err := ctx.SaveUploadedFile(body.ImgUrl, filepath); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})
	}

	savedPath := fmt.Sprintf("img/%s", filename)
	body.ImgUrl.Filename = savedPath

	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	if err := u.us.ChangeUserProfile(ctx.Request.Context(), body, claims.Id); err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.NoRowsAffected) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
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
