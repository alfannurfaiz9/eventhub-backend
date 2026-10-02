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

func initUserRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	r := router.Group("/user")

	am := middleware.NewAuthMiddleWare(rdb)

	ur := repo.NewUserRepo(db)
	us := service.NewUserService(ur)
	uc := controller.NewUserController(us)

	r.Use(am.CheckToken)

	r.GET("profile", am.UserMiddleware, uc.GetUserProfile)
	r.GET("event", am.UserMiddleware, uc.GetMyEvent)
	r.GET("notification", uc.GetNotification)
	r.POST("change-password", uc.ChangeUserPassword)
}
