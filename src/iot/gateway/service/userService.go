package service

import (
	"iot/gateway/model"
	"iot/gateway/logger"
	"iot/gateway/db"
)

type UserService struct {}

type UserInterface interface {
	Login(mobile string, password string) error
	UserDetail(userId int64) *model.User
}

// 登录
func (userService *UserService) Login(mobile string, password string) error  {
	logger := logger.GetLoggerInstance()
	var user = model.User{
		Mobile: mobile,
	}
	isExt, err := db.StdMasterDB().Get(&user)
	if err != nil {
		logger.Infof("query user info error, mobile: %s, error: %s", mobile, err.Error())
		return err
	}

	if isExt == true {
		logger.Infof("login successuful, mobile: %s", mobile)
		return err;
	}

	affected, err := db.StdMasterDB().Insert(&user)
	if err != nil {
		logger.Infof("lgoin info, mobile: %s, password: %s, error: %s", mobile, password, err.Error())
		return err
	}

	logger.Infof("login successuful, affected: %d", affected)

	return nil
}

// 用户详情
func (userService *UserService) UserDetail(userId int64) *model.User {
	var user = model.User{Id: 1}
	db.StdMasterDB().Get(user)

	return &user;
}


