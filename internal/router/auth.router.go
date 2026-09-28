package router

import (
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/controller"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/middleware"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	r := router.Group("/auth")

	ar := repo.NewAuthRepo(db)
	as := service.NewAuthService(ar)
	ac := controller.NewAuthController(as)

	r.POST("register", ac.Register)
	r.POST("login", ac.Login)
	r.POST("change-password", middleware.CheckToken, ac.ChangeUserPassword)
}
