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

func initAdminRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	r := router.Group("/admin/dashboard")

	am := middleware.NewAuthMiddleWare(rdb)

	ar := repo.NewAdminRepo(db)
	as := service.NewAdminService(ar)
	ac := controller.NewAdminController(as)

	r.Use(am.CheckToken, am.AdminMiddleware)

	r.GET("", ac.GetAdminDashboard)
	r.GET("users", ac.GetAllUser)
}
