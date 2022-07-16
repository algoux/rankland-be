package middleware

import (
	"net/http"
	"ranklist/errcode"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		err := ctx.Errors.Last()
		if err == nil {
			return
		}

		meta, ok := err.Meta.(errcode.Err)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "未知错误",
			})
			return
		}

		h := gin.H{
			"code":    meta.Code,
			"message": meta.Message,
		}
		switch err.Type {
		case gin.ErrorTypeBind:
			ctx.JSON(http.StatusBadRequest, h)
		case gin.ErrorTypeRender:
			ctx.JSON(http.StatusInternalServerError, h)
		default:
			ctx.JSON(http.StatusOK, h)
		}
	}
}
