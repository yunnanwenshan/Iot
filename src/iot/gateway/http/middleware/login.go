package middleware

import (
	"github.com/gin-gonic/gin"

	"iot/gateway/logger"
	"iot/gateway/util/serverToken"

	"net/http"
)

var (
	//clientInfoKey = "clientInfo"
	serverTokenKey = "Server-Token"
)

func NeedLogin() gin.HandlerFunc {
	return needLogin();
}

func needLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logger.GetLoggerInstance()
		serverTokenStr := ctx.Request.Header.Get(serverTokenKey)
		logger.Infof("Server-Token: %s", serverTokenStr)
		if ok, err := serverToken.VlidateServerToken(serverTokenStr); !ok {
			logger.Infof("ok=%v, err=%v", ok, err)
			ctx.JSON(http.StatusOK, gin.H{"code": 10008, "msg": "Server-token验证失败"})
			ctx.Abort()
		}

		ctx.Next()
	}
}
