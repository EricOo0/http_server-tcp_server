package impl

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"tcp_server/src/dao"
	"tcp_server/src/db"
	"tcp_server/src/logger"
	"tcp_server/src/model"
)

func init() {
	daoImpl := &IUserInfoDao{}
	dao.InjectUserInfoInstance(daoImpl)
}

type IUserInfoDao struct {
}

func (impl *IUserInfoDao) GetUserInfoByUserId(userId string) (*model.UserInfo, error) {
	users := make([]*model.UserInfo, 0)
	suffix, err := strconv.Atoi(userId)
	if err != nil {
		logger.DefaultLogger.Error("UserId invalid")
		return nil, errors.New("UserId invalid")
	}
	suffix = suffix % 10
	realTableName := "user_info_tab_0000000" + strconv.Itoa(suffix%10)
	session := db.GetDbInstance().NewSession()
	session.Table(realTableName).Where("user_id=?", userId).Find(&users)
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, nil
}
func (impl *IUserInfoDao) UpdateUserInfo(user model.UserInfoRecord, oldpassword string) int64 {
	suffix, err := strconv.Atoi(*user.UserId)
	if err != nil {
		logger.DefaultLogger.Error("UserId invalid")
		return 0
	}
	suffix = suffix % 10
	realTableName := "user_info_tab_0000000" + strconv.Itoa(suffix%10)
	session := db.GetDbInstance().NewSession()
	count, err := session.Table(realTableName).Where("user_id=? and password=?", *user.UserId, oldpassword).Update(user)
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("update userId %v error,err : %v", *user.UserId, err))
		return 0
	}
	return count
}

func (impl *IUserInfoDao) InsertUser(user model.UserInfo) int64 {
	suffix, err := strconv.Atoi(user.UserId)
	if err != nil {
		logger.DefaultLogger.Error("UserId invalid")
		return 0
	}
	suffix = suffix % 10
	realTableName := "user_info_tab_0000000" + strconv.Itoa(suffix%10)
	session := db.GetDbInstance().NewSession()
	count, err := session.Table(realTableName).Insert(user)
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("insert userId %v error,err : %v", user.UserId, err))
		return 0
	}
	return count
}
