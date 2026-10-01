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

// GetAdminDashboard
//
// @Summary			Get admin dashboard
// @Description		Get admin dashboard
// @Tags			admin
// @Produce			json
// @Router			/admin/dashboard	[get]
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (a *AdminController) GetAdminDashboard(ctx *gin.Context) {
	result, err := a.as.GetAdminDashboard(ctx.Request.Context())

	if err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: true,
			Message: "internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    result,
		Message: "successfully get admin dashboard",
	})
}

// GetAdminDashboardUser
//
// @Summary			Get admin dashboard all users
// @Description		Get admin dashboard all users
// @Tags			admin
// @Produce			json
// @Router			/admin/dashboard/users	[get]
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (a *AdminController) GetAllUser(ctx *gin.Context) {
	result, err := a.as.GetAllUser(ctx.Request.Context())

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
		Message: "successfully get all users",
	})
}
