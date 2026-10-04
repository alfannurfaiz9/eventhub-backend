package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

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
//	@Param        	location    query     string  false  "name location"  Format(location)
//	@Param        	category    query     string  false  "name category"  Format(category)
//	@Param        	page    query     int  false  "name page"  Format(page)
//
// @Router			/events	[get]
// @Success			200		{object}	dto.Response
// @Failure			404		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) GetEvents(ctx *gin.Context) {
	search := ctx.Query("search")
	location := ctx.Query("location")
	category := ctx.Query("category")
	page := ctx.DefaultQuery("page", "1")
	pageNum, _ := strconv.Atoi(page)

	events, err := e.es.GetEvents(ctx.Request.Context(), search, location, category, pageNum)

	if err != nil {
		if errors.Is(err, custom_error.EventErrorPage) {
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
// @Router			/events/{id}/join		[post]
// @Param        	id   path      int  true  "Event ID"
// @Security		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			404		{object}	dto.ErrorResponse
// @Failure			406		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) JoinEvent(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)
	eventId := ctx.Param("event_id")
	eIdInt, _ := strconv.Atoi(eventId)

	err := e.es.JoinEvent(ctx.Request.Context(), claims.Id, eIdInt)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.NoRowsAffected) {
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
// @Router			/events/{id}/leave		[delete]
// @Param        	id   path      int  true  "Event ID"
// @Security		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			406		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) LeaveEvent(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)
	eventId := ctx.Param("event_id")
	eIdInt, _ := strconv.Atoi(eventId)

	if err := e.es.LeaveEvent(ctx.Request.Context(), claims.Id, eIdInt); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "succesfully leave event",
	})
}

// SaveEvent
//
// @Summary			Save event
// @Description		Save specified event
// @Tags			events
// @Produce			json
// @Router			/events/{id}/save		[post]
// @Param        	id   path      int  true  "Event ID"
// @Security		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			406		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) SaveEvent(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)
	eventId := ctx.Param("event_id")
	eIdInt, _ := strconv.Atoi(eventId)

	if err := e.es.SaveEvent(ctx.Request.Context(), claims.Id, eIdInt); err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.NoRowsAffected) {
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
		Message: "successfully save event",
	})
}

// CreateEvent
//
// @Summary			Create new event
// @Description		Create new event
// @Tags			events
// @Accept			mpfd
// @Produce			json
// @Security		BearerToken
// @Router			/events/create	[post]
// @Param			title				formData	string		true	"title"
// @Param			bio					formData	string		true	"bio"
// @Param			img_url				formData	file		true	"img_url"
// @Param			description			formData	string		true	"description"
// @Param			start_at			formData	string		true	"start_at" 			format(date-time)
// @Param			end_at				formData	string		true	"end_at" 			format(date-time)
// @Param			format				formData	string		true	"format"
// @Param			capacity			formData	string		true	"capacity"
// @Param			community_id		formData	string		false	"community_id"
// @Param			category_id			formData	[]int		true	"category_id"		collectionFormat(multi)
// @Param			location_name		formData	string		true	"location_name"
// @Param			speaker_name		formData	[]string	true	"speaker_name" 		collectionFormat(multi)
// @Param			speaker_img_url		formData	[]string	false	"speaker_img_url" 	collectionFormat(multi)
// @Param			speaker_position	formData	[]string	true	"speaker_position"	collectionFormat(multi)
// @Param			speaker_company		formData	[]string	true	"speaker_company" 	collectionFormat(multi)
// @Success			201		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (e *EventController) CreateEvent(ctx *gin.Context) {
	var body dto.CreateEvent

	if err := ctx.ShouldBindWith(&body, binding.FormMultipart); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	var imgUrl string
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), "event", path.Ext(body.Img.Filename))
	filepath := path.Join("public", "img", filename)

	if err := ctx.SaveUploadedFile(&body.Img, filepath); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	imgUrl = filename

	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	if err := e.es.CreateEvent(ctx.Request.Context(), body, claims.Id, imgUrl); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "event successfully created",
	})
}
