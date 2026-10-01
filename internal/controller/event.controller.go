package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
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

func (e *EventController) GetEvents(ctx *gin.Context) {
	search := ctx.Query("search")
	location := ctx.Query("location")
	category := ctx.Query("category")

	events, err := e.es.GetEvents(ctx.Request.Context(), search, location, category)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    events,
		Message: "Success",
	})
}

func (e *EventController) GetEventDetail(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	result, err := e.es.GetEventDetail(ctx.Request.Context(), id)

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
		Message: "Success",
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
	id, _ := ctx.Get("id")

	var body dto.JoinEvent
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	err := e.es.JoinEvent(ctx.Request.Context(), id.(int), body)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "successfully join event",
	})
}

func (e *EventController) GetUpcomingEvent(ctx *gin.Context) {
	events, err := e.es.GetUpcomingEvent(ctx.Request.Context())

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
		Data:    events,
		Message: "Success",
	})
}

func (e *EventController) LeaveEvent(ctx *gin.Context) {
	userId, _ := ctx.Get("id")

	var body dto.JoinEvent
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := e.es.LeaveEvent(ctx.Request.Context(), userId.(int), body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Success",
	})
}
