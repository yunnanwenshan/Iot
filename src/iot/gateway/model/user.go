package model

import "time"

type User struct {
	Id uint64 `xorm:"pk autoincr"`
	Mobile string `xorm:"mobile"`
	Token string `xorm:"token"`
	CreateTime time.Time `xorm:"create_time"`
	UpdateTime time.Time `xorm:"update_time"`
}

func (user *User) TableName() string {
	return "user"
} 

type UserInfo struct {
	Id uint64
	UserId int64 `xorm:"user_id"`
}

func (userInfo *UserInfo) TableName() string {
	return "user_info"
}

type UserBind struct {
	UserName string `form:"user_name" json:"user_name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	DeviceSn string `form:"device_sn" json:"device_sn" binding:"required"`
}