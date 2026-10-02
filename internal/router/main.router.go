package router

import (
	_ "github.com/alfannurfaiz9/eventhub-backend.git/docs"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitMainRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	router.GET("documentation/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	initAuthRouter(router, db, rdb)
	initEventRouter(router, db, rdb)
	initCommunityRouter(router, db, rdb)
	initUserRouter(router, db, rdb)
	initTestimonyRouter(router, db, rdb)
	initOrganizerRouter(router, db, rdb)
	initAdminRouter(router, db, rdb)
}
