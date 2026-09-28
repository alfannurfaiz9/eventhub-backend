package router

import (
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/controller"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/middleware"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/repo"
	"github.com/alfannurfaiz9/eventhub-backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initTestimonyRouter(router *gin.Engine, db *pgxpool.Pool) {
	r := router.Group("/testimonies")

	tr := repo.NewTestimonyRepo(db)
	ts := service.NewTestimonyService(tr)
	tc := controller.NewTestimonyController(ts)

	r.GET("", tc.GetTestimony)
	r.POST("", middleware.CheckToken, tc.SetTestimony)
}
