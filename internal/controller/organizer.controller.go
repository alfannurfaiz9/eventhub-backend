package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
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

func (o *OrganizerController) GetOrganizerDashboard(ctx *gin.Context) {
	id, _ := ctx.Get("id")
	result, err := o.os.GetOrganizerDashboard(ctx.Request.Context(), id.(int))

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
		Message: "succeffully get dashboard",
	})
}

func (o *OrganizerController) GetOrganizerEvent(ctx *gin.Context) {
	id, _ := ctx.Get("id")
	result, err := o.os.GetOrganizerEvent(ctx.Request.Context(), id.(int))

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.EventNotFound) {
			ctx.JSON(http.StatusForbidden, dto.Response{
				Success: false,
				Message: err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: true,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get event",
	})
}
