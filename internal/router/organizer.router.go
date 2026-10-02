package router

import (
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/controller"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/middleware"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func initOrganizerRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	r := router.Group("/organizer/dashboard")

	am := middleware.NewAuthMiddleWare(rdb)

	or := repo.NewOrganizerRepo(db)
	os := service.NewOrganizerService(or)
	oc := controller.NewOrganizerController(os)

	r.Use(am.UserMiddleware)

	r.GET("", oc.GetOrganizerDashboard)
	r.GET("events", oc.GetOrganizerEvent)
}
