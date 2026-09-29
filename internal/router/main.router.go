package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitMainRouter(router *gin.Engine, db *pgxpool.Pool) {
	initAuthRouter(router, db)
	initEventRouter(router, db)
	initCommunityRouter(router, db)
	initUserRouter(router, db)
	initTestimonyRouter(router, db)
	initOrganizerRouter(router, db)
}
