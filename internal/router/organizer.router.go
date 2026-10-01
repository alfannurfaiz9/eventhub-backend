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

	or := repo.NewOrganizerRepo(db)
	os := service.NewOrganizerService(or)
	oc := controller.NewOrganizerController(os)

	r.Use(middleware.CheckToken)

	r.GET("dashboard", oc.GetOrganizerDashboard)
	r.GET("events", oc.GetOrganizerEvent)
}
