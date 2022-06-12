package dao

import (
	"tcp_server/src/model"
)

var userInfoDaoInstance UserInfoDao

type UserInfoDao interface {
	GetUserInfoByUserId(userId string) (*model.UserInfo, error)
	UpdateUserInfo(user model.UserInfoRecord, oldPassword string) int64
	InsertUser(user model.UserInfo) int64
}

func InjectUserInfoInstance(_userInfoDao UserInfoDao) {
	userInfoDaoInstance = _userInfoDao
}
func GetUserInfoDaoInstance() UserInfoDao {
	return userInfoDaoInstance
}
