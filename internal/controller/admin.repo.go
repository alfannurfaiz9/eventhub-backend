package controller

import (
	"log"
	"net/http"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/dto"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	as *service.AdminService
}

func NewAdminController(as *service.AdminService) *AdminController {
	return &AdminController{
		as: as,
	}
}

func (a *AdminController) GetAdminDashboard(ctx *gin.Context) {
	result, err := a.as.GetAdminDashboard(ctx.Request.Context())

	if err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: true,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get admin dashboard",
	})
}

func (a *AdminController) GetAllUser(ctx *gin.Context) {
	result, err := a.as.GetAllUser(ctx.Request.Context())

	if err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Data:    result,
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Success: true,
		Data:    result,
		Message: "successfully get all users",
	})
}
