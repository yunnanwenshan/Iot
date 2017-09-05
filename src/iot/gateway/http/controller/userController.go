package controller

import (
	"github.com/gin-gonic/gin"
	"iot/gateway/logger"
	"iot/gateway/util/serverToken"
	//"api/http/middleware"
	"iot/gateway/service"
	"iot/gateway/model"
	"strconv"
)

type UserController struct {}

var (
	userService = new(service.UserService)
)

func (self UserController) RegisterRouter(r *gin.Engine)  {
	//需要登录才可以访问的接口
	//loginRouter := r.Group("app/v1/user/", middleware.NeedLogin())
	loginRouter := r.Group("app/v1/user/")
	loginRouter.POST("bind", self.BindDevice)
	//loginRouter.POST("/user/:user_id", self.UserDetail)

	//不登录也可以访问的接口
	//g := r.Group("app/v1/user")
	//g.POST("/login", self.Login)
}

// 用户详情
func (self UserController) UserDetail(ctx *gin.Context) {
	logger := logger.GetLoggerInstance()
	userId := ctx.Param("user_id")
	logger.Infof("userId: %d, test===============", userId)
	userDetail := userService.UserDetail(1234)
	success(ctx, userDetail)
}

// 登录接口
func (self UserController) Login(ctx *gin.Context) {
	logger := logger.GetLoggerInstance()
	token, err := serverToken.GenerateToken(string(98))

	if err != nil {
		logger.Warnf("生成token失败, error: %v", err)
		fail(ctx, 1000, "登录失败")
		return
	}

	//userService := new(service.UserService)
	err = userService.Login("13311588124", "1234")
	if err != nil {
		logger.Infof("login fail, errr: %s", err.Error())
		fail(ctx, 2000, "登录失败")
		return
	}

	success(ctx, map[string]interface{}{"ticket":token})
}

//绑定设备
func (self UserController) BindDevice(ctx *gin.Context) {
	var userBind model.UserBind
	logger := logger.GetLoggerInstance()
	logger.Infof("bind device begin.......")
	err := ctx.BindJSON(&userBind)
	if err != nil {
		logger.Infof("get params fail, err: %s", err.Error())
		fail(ctx, 1002, "获取参数失败")
		return
	}
	deviceNum , err := strconv.Atoi(userBind.DeviceSn)
	if err != nil {
		logger.Infof("get params fail 1, err: %s", err.Error())
		fail(ctx, 1003, "获取参数失败")
		return
	}
	user, err := userService.BindUserToDevice(userBind.UserName, userBind.Password, deviceNum)
	if err != nil {
		logger.Errorf("bind device fail, err: %s", err.Error())
		fail(ctx, 1001, "绑定失败")
		return
	}

	success(ctx, map[string]interface{}{"ticket": user.Token, "user_id": user.Id, "mobile": user.Mobile})
}
