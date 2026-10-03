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

func initCommunityRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	r := router.Group("/communities")

	am := middleware.NewAuthMiddleWare(rdb)

	cr := repo.NewCommunityRepo(db)
	cs := service.NewCommunityService(cr)
	cc := controller.NewCommunityController(cs)

	r.GET("", cc.GetCommunities)
	r.GET(":id", cc.GetCommunityDetail)
	r.GET(":id/events", cc.GetCommunityEvent)
	r.GET(":id/members", cc.GetCommunityMember)
	r.GET("popular", cc.GetPopularCommunity)

	r.POST(":community_id/join", am.CheckToken, am.UserMiddleware, cc.JoinCommunity)
	r.DELETE(":community_id/leave", am.CheckToken, am.UserMiddleware, cc.LeaveCommunity)
}
