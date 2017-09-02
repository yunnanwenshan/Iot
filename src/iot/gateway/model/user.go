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