package router

import (
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/controller"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/middleware"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initAdminRouter(router *gin.Engine, db *pgxpool.Pool) {
	r := router.Group("/admin/dashboard")

	ar := repo.NewAdminRepo(db)
	as := service.NewAdminService(ar)
	ac := controller.NewAdminController(as)

	r.GET("", middleware.CheckToken, ac.GetAdminDashboard)
	r.GET("users", middleware.CheckToken, ac.GetAllUser)
}
