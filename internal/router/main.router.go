package router

import (
	_ "github.com/alfannurfaiz9/eventhub-backend.git/docs"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitMainRouter(router *gin.Engine, db *pgxpool.Pool) {
	router.GET("documentation/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	initAuthRouter(router, db)
	initEventRouter(router, db)
	initCommunityRouter(router, db)
	initUserRouter(router, db)
	initTestimonyRouter(router, db)
	initOrganizerRouter(router, db)
	initAdminRouter(router, db)
}
