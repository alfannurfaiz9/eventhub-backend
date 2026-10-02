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

func initTestimonyRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	r := router.Group("/testimonies")

	am := middleware.NewAuthMiddleWare(rdb)

	tr := repo.NewTestimonyRepo(db)
	ts := service.NewTestimonyService(tr)
	tc := controller.NewTestimonyController(ts)

	r.GET("", tc.GetTestimony)
	r.POST("", am.UserMiddleware, tc.SetTestimony)
}
