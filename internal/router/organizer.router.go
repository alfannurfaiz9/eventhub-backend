package router

import (
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/controller"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/middleware"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initOrganizerRouter(router *gin.Engine, db *pgxpool.Pool) {
	r := router.Group("/organizer")

	ar := repo.NewAuthRepo(db)
	as := service.NewAuthService(ar)
	ac := controller.NewAuthController(as)

	r.GET("dashboard", middleware.CheckToken, ac.GetOrganizerDashboard)
}
