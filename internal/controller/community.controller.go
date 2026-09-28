package controller

import (
	"log"
	"net/http"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/alfannurfaiz9/eventhub-backend.git/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CommunityController struct {
	cs *service.CommunityService
}

func NewCommunityController(cs *service.CommunityService) *CommunityController {
	return &CommunityController{
		cs: cs,
	}
}

func (c *CommunityController) GetCommunities(ctx *gin.Context) {
	category := ctx.Query("category")
	result, err := c.cs.GetCommunities(ctx.Request.Context(), category)

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
		Message: "Success",
	})
}

func (c *CommunityController) GetCommunityDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.cs.GetCommunityDetail(ctx.Request.Context(), id)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusInternalServerError, dto.Response{
		Success: true,
		Data:    result,
		Message: "success",
	})
}

func (c *CommunityController) GetCommunityEvent(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.cs.GetCommunityEvent(ctx, id)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusInternalServerError, dto.Response{
		Success: true,
		Data:    result,
		Message: "success",
	})
}

func (c *CommunityController) GetCommunityMember(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.cs.GetCommunityMember(ctx, id)

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

func (c *CommunityController) GetPopularCommunity(ctx *gin.Context) {
	result, err := c.cs.GetPopularCommunity(ctx)

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

func (c *CommunityController) JoinCommunity(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	claims, _ := token.(pkg.JWTClaims)

	var body dto.UserCommunity
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	err := c.cs.JoinCommunity(ctx.Request.Context(), claims.Id, body)

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
		Message: "successfully join community",
	})
}

func (c *CommunityController) LeaveCommunity(ctx *gin.Context) {
	userId, _ := ctx.Get("id")

	var body dto.UserCommunity
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := c.cs.LeaveCommunity(ctx.Request.Context(), userId.(int), body); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "successfully leave community",
	})
}
