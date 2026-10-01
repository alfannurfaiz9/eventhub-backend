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

func initEventRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	r := router.Group("/events")

	er := repo.NewEventRepo(db)
	es := service.NewEventService(er, rdb)
	ec := controller.NewEventController(es)

	r.GET("", ec.GetEvents)
	r.GET("upcoming", ec.GetUpcomingEvent)
	r.GET("detail/:id", ec.GetEventDetail)
	r.POST("join", middleware.CheckToken, ec.JoinEvent)
	r.DELETE("leave", middleware.CheckToken, ec.LeaveEvent)
}
