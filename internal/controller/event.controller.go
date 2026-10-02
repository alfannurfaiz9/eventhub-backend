package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	custom_error "github.com/alfannurfaiz9/eventhub-backend.git/internal/error"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type EventController struct {
	es *service.EventService
}

func NewEventController(es *service.EventService) *EventController {
	return &EventController{
		es: es,
	}
}

// GetEvents
//
// @Summary			Get all event
// @Description		Get all event with search and filter
// @Tags			events
// @Produce			json
//
//	@Param        	search    query     string  false  "name search by search"  Format(search)
//	@Param        	location    query     string  false  "name location by location"  Format(location)
//	@Param        	category    query     string  false  "name category by category"  Format(category)
//
// @Router			/events	[get]
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) GetEvents(ctx *gin.Context) {
	search := ctx.Query("search")
	location := ctx.Query("location")
	category := ctx.Query("category")

	events, err := e.es.GetEvents(ctx.Request.Context(), search, location, category)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    events,
		Message: "successfully get event",
	})
}

// GetEventDetail
//
// @Summary			Get event detail
// @Description		Get event detail
// @Tags			events
// @Produce			json
// @Router			/events/detail/{id}	[get]
// @Param        	id   path      int  true  "Event ID"
// @Success			200		{object}	dto.Response
// @Failure			404		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) GetEventDetail(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	result, err := e.es.GetEventDetail(ctx.Request.Context(), id)

	if err != nil {
		log.Println(err.Error())

		if errors.Is(err, custom_error.EventNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
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
		Message: "successfully get event detail",
	})
}

// JoinEvent
//
// @Summary			Join to event
// @Description		Join to specified event
// @Tags			events
// @Produce			json
// @Router			/events/join		[post]
// @Param			data	body		dto.JoinEvent	true	"body to join event"
// @Security		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) JoinEvent(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	var body dto.JoinEvent
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	err := e.es.JoinEvent(ctx.Request.Context(), claims.Id, body)

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
		Message: "successfully join event",
	})
}

// UpcomingEvent
//
// @Summary			Get upcoming event
// @Description		Get upcoming event
// @Tags			events
// @Produce			json
// @Router			/events/upcoming	[get]
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) GetUpcomingEvent(ctx *gin.Context) {
	events, err := e.es.GetUpcomingEvent(ctx.Request.Context())

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
		Data:    events,
		Message: "successfully get upcoming event",
	})
}

// LeaveEvent
//
// @Summary			Leave event
// @Description		Leave specified event
// @Tags			events
// @Produce			json
// @Router			/events/leave		[delete]
// @Param			data	body		dto.JoinEvent	true	"body to leave event"
// @Security		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) LeaveEvent(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	var body dto.JoinEvent
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	if err := e.es.LeaveEvent(ctx.Request.Context(), claims.Id, body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Success",
	})
}
