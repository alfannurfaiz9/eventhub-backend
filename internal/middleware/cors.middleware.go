package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	ctx.Header("Access-Control-Allow-Methods", "OPTIONS, POST, DELETE, PATCH, PUT")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
