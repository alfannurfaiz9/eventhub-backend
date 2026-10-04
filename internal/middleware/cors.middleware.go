package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "http://localhost:5500")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, XXX-Header")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
