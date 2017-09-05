package service

import (
	"iot/gateway/model"
	"iot/gateway/logger"
	"iot/gateway/db"
	"iot/gateway/redis"
	"iot/internal/rds"
	"strconv"
	"math/rand"
	"github.com/pkg/errors"
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
		logger.Infof("login info, mobile: %s, password: %s, error: %s", mobile, password, err.Error())
		return err
	}

	logger.Infof("login successuful, affected: %d", affected)

	return nil
}

// 用户详情
func (userService *UserService) UserDetail(userId uint64) *model.User {
	user := model.User{
		Id: userId,
	}
	db.StdMasterDB().Get(user)

	return &user;
}

//根据手机号查询用户信息
func (userService *UserService) UserDetailByUserName(userName string) *model.User {
	user := model.User{
		Mobile: userName,
	}
	db.StdMasterDB().Get(user)
	return &user
}

//设备与帐号绑定
func (userService *UserService) BindUserToDevice(userName string, password string, deviceNum int) (*model.User, error) {
	logger := logger.GetLoggerInstance()
	user := new(model.User)
	user.Id = 1001;
	uid := strconv.Itoa(int(user.Id))
	token := strconv.Itoa(rand.Int())
	user.Token = token
	user.Mobile = userName


	redisInstance := redis.RedisForDeviceInstance()
	//查找设备是否已经绑定
	ss, err := redisInstance.FindSessions(uid)
	if err != nil {
		logger.Warnf("read user info fail from redis, err: %s", err.Error())
		//建立设备与用户的关心
		session := rds.Session{
			Id: uid,
			AuthCode: token,
			Login: true,
			Online: true,
			Plat: deviceNum,
		}
		sessiones := rds.Sessions{
			Id: uid,
			Sess: []*rds.Session{&session},
		}
		redisInstance.SaveSessions(&sessiones)
	} else {
		session := ss.Sess[0]
		session.AuthCode = token
		if session.Plat != deviceNum {
			logger.Infof("the user binded to other device")
			return nil, errors.New("the user binded to other device")
		}
		redisInstance.SaveSessions(ss)
	}

	logger.Infof("bind device successful, user_id: %d, token: %s", user.Id, token)

	return user, nil
}


