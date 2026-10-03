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
)

type CommunityController struct {
	cs *service.CommunityService
}

func NewCommunityController(cs *service.CommunityService) *CommunityController {
	return &CommunityController{
		cs: cs,
	}
}

// GetCommunities
//
// @Summary			Get all community
// @Description		Get all community with search and filter
// @Tags			communities
// @Produce			json
//
//	@Param        	category    query     string  false  "name category by category"  Format(category)
//	@Param        	page    query     int  false  "name page"  Format(page)
//
// @Router			/communities	[get]
// @Success			200		{object}	dto.Response
// @Success			404		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityController) GetCommunities(ctx *gin.Context) {
	category := ctx.Query("category")
	page := ctx.DefaultQuery("page", "1")
	pageNum, _ := strconv.Atoi(page)
	result, err := c.cs.GetCommunities(ctx.Request.Context(), category, pageNum)

	if err != nil {
		log.Println(err.Error())
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
		Data:    result,
		Message: "Success",
	})
}

// GetCommunityDetail
//
// @Summary			Get community detail
// @Description		Get community detail
// @Tags			communities
// @Produce			json
// @Router			/communities/{id}	[get]
// @Param        	id   path      int  true  "Community ID"
// @Success			200		{object}	dto.Response
// @Failure			404		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityController) GetCommunityDetail(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	result, err := c.cs.GetCommunityDetail(ctx.Request.Context(), id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.CommunityNotFound) {
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
		Message: "successfully get community detail",
	})
}

// GetCommunityEvent
//
// @Summary			Get community event
// @Description		Get community event
// @Tags			communities
// @Produce			json
// @Router			/communities/{id}/events	[get]
// @Param        	id   path      int  true  "Community ID"
// @Success			200		{object}	dto.Response
// @Failure			404		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityController) GetCommunityEvent(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	result, err := c.cs.GetCommunityEvent(ctx.Request.Context(), id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, custom_error.CommunityNotFound) {
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
		Message: "successfully get event",
	})
}

// GetCommunityMember
//
// @Summary			Get community member
// @Description		Get community member
// @Tags			communities
// @Produce			json
// @Param        	id   path      int  true  "Community ID"
// @Router			/communities/{id}/members	[get]
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityController) GetCommunityMember(ctx *gin.Context) {
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	result, err := c.cs.GetCommunityMember(ctx.Request.Context(), id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get community member",
	})
}

// GetPopularCommunity
//
// @Summary			Get popular community
// @Description		Get popular community
// @Tags			communities
// @Produce			json
// @Router			/communities/popular	[get]
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityController) GetPopularCommunity(ctx *gin.Context) {
	result, err := c.cs.GetPopularCommunity(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get popular community",
	})
}

// JoinCommunity
//
// @Summary			Join to community
// @Description		Join to specified community
// @Tags			communities
// @Produce			json
// @Router			/communities/{id}/join		[post]
// @Param        	id   path      int  true  "Community ID"
// @Security		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityController) JoinCommunity(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)
	communitytId := ctx.Param("community_id")
	cIdInt, _ := strconv.Atoi(communitytId)

	err := c.cs.JoinCommunity(ctx.Request.Context(), claims.Id, cIdInt)

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
		Message: "successfully join community",
	})
}

// LeaveCommunity
//
// @Summary			Leave community
// @Description		Leave specified community
// @Tags			communities
// @Produce			json
// @Router			/communities/{id}/leave		[delete]
// @Param        	id   path      int  true  "Community ID"
// @Security		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityController) LeaveCommunity(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)
	communitytId := ctx.Param("community_id")
	cIdInt, _ := strconv.Atoi(communitytId)

	if err := c.cs.LeaveCommunity(ctx.Request.Context(), claims.Id, cIdInt); err != nil {
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
		Message: "successfully leave community",
	})
}
