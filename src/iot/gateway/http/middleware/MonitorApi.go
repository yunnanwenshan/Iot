package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
	"iot/gateway/logger"
)

func MonitorApi() gin.HandlerFunc {
	return statTime()
}

func statTime() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now().UnixNano()
		logger := logger.GetLoggerInstance()

		//继续处理下一个中间件
		ctx.Next()

		//计算接口所花费的时间
		end := time.Now().UnixNano()
		diff := (end - begin) / (1000 * 1000)

		if ((diff / 1000) >= 1) {
			logger.Warnf("uri: %s, warning public cost time: %v ms", ctx.Request.URL.Path, (diff / 1000))
		} else {
			logger.Infof("uri: %s, public cost time: %v ms", ctx.Request.URL.Path, diff)
		}
	}
}
