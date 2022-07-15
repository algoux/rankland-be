package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 允许跨域请求
func Cors(cors string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", cors)
		ctx.Writer.Header().Set("Access-Control-Max-age", "3600")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusMethodNotAllowed)
		}
	}
}
