package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"strconv"
	"tcp_server/src/cache"
	"tcp_server/src/dao"
	"tcp_server/src/logger"
	"tcp_server/src/model"
	api "tcp_server/src/proto/proto"
	"tcp_server/src/server_config"
	"time"
)

var tcpServer *grpc.Server

func InitTcpServer() {
	tcpServer = grpc.NewServer()
	api.RegisterAdminServiceServer(tcpServer, &TcpServer{})
	logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server Init successfully,listen on %v:%v", server_config.GlbserverConfig.TcpHost, server_config.GlbserverConfig.TcpPort))
}

func GetTcpServer() *grpc.Server {
	return tcpServer
}

type TcpServer struct {
}

func (s *TcpServer) UserInfo(ctx context.Context, req *api.UserInfoReq) (*api.UserInfoRsp, error) {
	logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -GetUserInfo"))
	defer func() {
		if rec := recover(); rec != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("err:%v", rec))
		}
		logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -GetUserInfo done"))

	}()
	user, err := dao.GetUserInfoDaoInstance().GetUserInfoByUserId(req.UserId)
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("Get UserInfo failed,err:%v", err))
		return nil, err
	}
	logger.DefaultLogger.Info(fmt.Sprintf("Get UserInfo success,id:%v,nickname:%v", user.UserId, user.Nickname))

	rsp := &api.UserInfoRsp{
		UserId:         user.UserId,
		NickName:       user.Nickname,
		ProfilePicture: user.ProfilePicture,
	}
	return rsp, nil
}

func (s *TcpServer) EditUserInfo(ctx context.Context, req *api.UserEditReq) (*api.UserEditRsp, error) {
	logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -EditUserInfo"))
	defer func() {
		if rec := recover(); rec != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("err:%v", rec))
		}
		logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -EditUserInfo done"))

	}()
	user := model.UserInfo{
		UserId:          req.UserId,
		Password:        req.NewPassword,
		Nickname:        req.NickName,
		UpdateTimestamp: time.Now().Unix(),
	}
	userRecord := model.GenerateUserinfoRecord(user)
	// 获取缓存的用户密码
	curPassword, err := cache.Rdb.Get(fmt.Sprintf("%v_password", req.UserId)).Result()
	if err == redis.Nil {
		//不在缓存去db查
		logger.DefaultLogger.Info(fmt.Sprintf("userId '%v' 's password is not exist in cache,check in db", req.UserId))
		user, err := dao.GetUserInfoDaoInstance().GetUserInfoByUserId(req.UserId)
		if err != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("GetUserInfo error:%v", err))
			return &api.UserEditRsp{Status: false, Msg: fmt.Sprintf("GetUserInfo error:%v", err)}, nil
		}
		if user == nil {
			logger.DefaultLogger.Error("User not exist")
			return &api.UserEditRsp{Status: false, Msg: "User not exist"}, nil
		}
		curPassword = user.Password
	} else if err != nil {
		// 获取出错
		logger.DefaultLogger.Error(fmt.Sprintf("EditUserInfo Error ,err:%v", err))
		return &api.UserEditRsp{Status: false, Msg: fmt.Sprintf("%v", err)}, nil
	}

	if curPassword == req.OldPassword {
		count := dao.GetUserInfoDaoInstance().UpdateUserInfo(userRecord, curPassword)
		result := count != 0

		//删除缓存
		logger.DefaultLogger.Info(fmt.Sprintf("Del password cache of user:%v", *userRecord.UserId))
		_, err := cache.Rdb.Del(fmt.Sprintf("%v_password", *userRecord.UserId)).Result()
		if err != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("Del password cache of user:%v failed", *userRecord.UserId, err))
		}

		return &api.UserEditRsp{Status: result}, nil
	}
	return &api.UserEditRsp{Status: false, Msg: "Old Password unmatched"}, nil
}
func (s *TcpServer) UploadProfile(ctx context.Context, req *api.UploadProfileReq) (*api.UploadProfileRsp, error) {
	logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -UploadProfile"))
	defer func() {
		if rec := recover(); rec != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("err:%v", rec))
		}
		logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -UploadProfile done"))

	}()
	updateTime := time.Now().Unix()
	userRecord := model.UserInfoRecord{
		UserId:          &req.UserId,
		ProfilePicture:  &req.Url,
		UpdateTimestamp: &updateTime,
	}
	// 避免别人修改了密码这边还能修改
	curPassword, err := cache.Rdb.Get(fmt.Sprintf("%v_password", req.UserId)).Result()
	if err != nil {
		//不在缓存 重新登录
		logger.DefaultLogger.Error(fmt.Sprintf("userId '%v' 's password is not exist in cache,login again", req.UserId))
		return &api.UploadProfileRsp{Status: false, Msg: "user info expire,login again"}, nil
	}
	count := dao.GetUserInfoDaoInstance().UpdateUserInfo(userRecord, curPassword)
	result := count != 0
	if result {
		logger.DefaultLogger.Info(fmt.Sprintf("userId '%v' upload picture sucessfully", req.UserId))
		return &api.UploadProfileRsp{Status: result}, nil
	}
	return &api.UploadProfileRsp{Status: result, Msg: "password has been change,Please login again and update"}, nil
}

func (s *TcpServer) UserLogin(ctx context.Context, req *api.UserLoginReq) (*api.UserLoginRsp, error) {
	logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -UserLogin"))
	defer func() {
		if rec := recover(); rec != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("err:%v", rec))
		}
		logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -UserLogin done"))

	}()
	curPassword, err := cache.Rdb.Get(fmt.Sprintf("%v_password", req.UserId)).Result()
	if err == redis.Nil {
		//不在缓存去db查
		logger.DefaultLogger.Info(fmt.Sprintf("userId:%v is not exist in cache,check in db", req.UserId))
		user, err := dao.GetUserInfoDaoInstance().GetUserInfoByUserId(req.UserId)
		if err != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("GetUserInfo error:%v", err))
		}
		if user == nil {
			logger.DefaultLogger.Error("User not exist")
			return &api.UserLoginRsp{Status: false, Msg: "User not exist"}, nil
		}
		curPassword = user.Password
		logger.DefaultLogger.Info(fmt.Sprintf("set password cache ：userId:%v", req.UserId))
		cache.Rdb.Set(fmt.Sprintf("%v_password", req.UserId), curPassword, 24*time.Hour)
	} else if err != nil {
		// 获取出错
		logger.DefaultLogger.Error(fmt.Sprintf("EditUserInfo Error ,err:%v", err))
		return &api.UserLoginRsp{Status: false, Msg: fmt.Sprintf("%v", err)}, nil
	}

	if curPassword != req.Password {
		logger.DefaultLogger.Error(fmt.Sprintf("userId:%v；password not match", req.UserId))
		return &api.UserLoginRsp{Status: false, Msg: "password not match"}, nil
	}
	cache.Rdb.Expire(fmt.Sprintf("%v_password", req.UserId), 24*time.Hour)
	return &api.UserLoginRsp{Status: true}, nil
}

func (s *TcpServer) RegisterUser(ctx context.Context, req *api.RegisterUserReq) (*api.RegisterUserRsp, error) {
	logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -RegisterUser"))
	defer func() {
		if rec := recover(); rec != nil {
			logger.DefaultLogger.Error(fmt.Sprintf("err:%v", rec))
		}
		logger.DefaultLogger.Info(fmt.Sprintf("Tcp Server -RegisterUser done"))
	}()
	rsp := &api.RegisterUserRsp{
		Status: false,
	}
	if _, err := strconv.ParseInt(req.UserId, 10, 64); err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("userId '%v' invalid,err:%v", req.UserId, err))
		rsp.Msg = "userId invalid"
		return rsp, nil
	}

	userRecord := model.UserInfo{
		UserId:          req.UserId,
		Password:        req.Password,
		Nickname:        req.NickName,
		CreateTimestamp: time.Now().Unix(),
		UpdateTimestamp: time.Now().Unix(),
	}

	user, err := dao.GetUserInfoDaoInstance().GetUserInfoByUserId(req.UserId)
	if err != nil {
		logger.DefaultLogger.Error(fmt.Sprintf("Check userinfo failed,err:%v", err))
		rsp.Msg = fmt.Sprintf("%v", err)
		return rsp, nil
	}
	if user != nil {
		logger.DefaultLogger.Error("user id exist")
		err = errors.New("user id exist")
		rsp.Msg = fmt.Sprintf("%v", err)
		return rsp, nil
	}
	count := dao.GetUserInfoDaoInstance().InsertUser(userRecord)
	rsp.Status = count != 0
	return rsp, nil
}
