package middleware

import (
	"encoding/xml"
	"strings"

	"github.com/gin-gonic/gin"
)

// XMLHeader
func XMLHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.HasSuffix(ctx.Request.URL.Path, ".xml") && ctx.NegotiateFormat(gin.MIMEXML) == gin.MIMEXML {
			ctx.Writer.Write([]byte(xml.Header))
		}

		ctx.Next()
	}
}
