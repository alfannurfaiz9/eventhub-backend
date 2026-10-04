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

func initAuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	r := router.Group("/auth")

	am := middleware.NewAuthMiddleWare(rdb)

	ar := repo.NewAuthRepo(db)
	as := service.NewAuthService(ar, rdb)
	ac := controller.NewAuthController(as)

	r.POST("register", ac.Register)
	r.POST("login", ac.Login)
	r.POST("logout", am.CheckToken, ac.Logout)
	r.PATCH("forgot-password", ac.ForgotPassword)
}
