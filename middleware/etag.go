package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FileEtag() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val := ctx.GetHeader("algoux")
		if val != "rankland" {
			ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		}
	}
}
