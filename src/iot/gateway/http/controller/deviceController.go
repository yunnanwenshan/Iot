package controller

import (
	"github.com/gin-gonic/gin"
	"iot/gateway/rpc"
	"iot/gateway/logger"
)

type DeviceController struct {}

func (self *DeviceController) RegisterRouter(r *gin.Engine) {
	deviceGroup := r.Group("/app/v1/device/")
	deviceGroup.POST("info", self.DeviceInfo)
}

//获取设备信息
func (self *DeviceController) DeviceInfo(ctx *gin.Context) {
	logger := logger.GetLoggerInstance()
	rpcClient := rpc.GetNodeApiInstance()
	nodeId, err := rpcClient.GetNode("1001")
	if err != nil {
		logger.Errorf("get node number failed, err: %s", err.Error())
		fail(ctx, 3001, "获取节点信息失败")
		return
	}

	success(ctx, map[string]interface{}{
		"node_id": nodeId,
	})
}

//发送命令到设备
func (self *DeviceController) SendCmd(ctx *gin.Context) {

}
