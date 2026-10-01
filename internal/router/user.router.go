package router

import (
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/controller"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/middleware"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initUserRouter(router *gin.Engine, db *pgxpool.Pool) {
	r := router.Group("/user")

	ur := repo.NewUserRepo(db)
	us := service.NewUserService(ur)
	uc := controller.NewUserController(us)

	r.GET("profile", middleware.CheckToken, uc.GetUserProfile)
	r.GET("event", middleware.CheckToken, uc.GetMyEvent)
	r.GET("notification", middleware.CheckToken, uc.GetNotification)
}
