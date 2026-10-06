package router

import (
	"path"

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

	am := middleware.NewAuthMiddleWare(rdb)

	er := repo.NewEventRepo()
	es := service.NewEventService(er, rdb, db)
	ec := controller.NewEventController(es)

	r.GET("", ec.GetEvents)
	r.GET("upcoming", ec.GetUpcomingEvent)
	r.GET("detail/:id", ec.GetEventDetail)
	r.Static("img", path.Join("public", "img"))
	r.POST(":event_id/join", am.CheckToken, am.UserMiddleware, ec.JoinEvent)
	r.DELETE(":event_id/leave", am.CheckToken, am.UserMiddleware, ec.LeaveEvent)
	r.POST(":event_id/save", am.CheckToken, am.UserMiddleware, ec.SaveEvent)

	r.POST("create", am.CheckToken, am.OrganizerMiddleware, ec.CreateEvent)
}
