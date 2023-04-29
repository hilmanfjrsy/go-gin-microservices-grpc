package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ResponseError(context *gin.Context, code int, message string) {
	context.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"error":   http.StatusText(code),
		"message": strings.Replace(message, "rpc error: code = Unknown desc = ", "", -1),
	})
}

func ResponseSuccess(context *gin.Context, code int, data interface{}) {
	context.JSON(code, data)
}
