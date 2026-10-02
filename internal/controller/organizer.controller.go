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
)

type OrganizerController struct {
	os *service.OrganizerService
}

func NewOrganizerController(os *service.OrganizerService) *OrganizerController {
	return &OrganizerController{
		os: os,
	}
}

// GetOrganizerDashboard
//
// @Summary			Get organizer dashboard
// @Description		Get organizer dashboard
// @Tags			organizer
// @Produce			json
// @Router			/organizer/dashboard	[get]
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (o *OrganizerController) GetOrganizerDashboard(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	result, err := o.os.GetOrganizerDashboard(ctx.Request.Context(), claims.Id)

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
		Message: "succeffully get dashboard",
	})
}

// GetOrganizerDashboardEvent
//
// @Summary			Get organizer dashboard all users
// @Description		Get organizer dashboard all users
// @Tags			organizer
// @Produce			json
// @Router			/organizer/dashboard/events	[get]
// @Success			200		{object}	dto.Response
// @Failure			403		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (o *OrganizerController) GetOrganizerEvent(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	result, err := o.os.GetOrganizerEvent(ctx.Request.Context(), claims.Id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.EventNotFound) {
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: true,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get event",
	})
}
