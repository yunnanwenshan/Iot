package routes

import (
	"github.com/gin-gonic/gin"
	"iot/gateway/http/controller"
)

func RegisterRouters(r *gin.Engine)  {
	//注册路由
	new(controller.UserController).RegisterRouter(r)
	new(controller.DeviceController).RegisterRouter(r)
}
