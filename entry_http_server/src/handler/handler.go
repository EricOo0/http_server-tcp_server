package handler

import (
	"github.com/go-redis/redis"
	"http_server/src/cache"
	constant "http_server/src/common"
	"http_server/src/logger"
	"http_server/src/middleware"
	api "http_server/src/proto/proto"
	"http_server/src/service"
	"time"

	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
)

var log *zap.Logger

type UserParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}
type UpdateUserMask struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	NickName    string `json:"nickname"`
}

func InitHandler() {
	log = logger.DefaultLogger
}

func Login(c *gin.Context) {
	log.Info("User Invoke Login API")
	defer func() {
		log.Info("User Finish  Login API")
	}()
	userInfo := UserParam{}
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		log.Error(fmt.Sprintf("phase login json failed,err:%v", err))
		c.JSON(constant.InvalidParam, gin.H{
			"msg": "failed,please check param",
		})
		return
	}
	log.Info(fmt.Sprintf("User '%v' Invloke TcpClient:UserLogin", userInfo.Username))
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(userInfo.Password))
	encryptPassword := hex.EncodeToString(md5Ctx.Sum(nil))
	req := &api.UserLoginReq{
		UserId:   userInfo.Username,
		Password: encryptPassword,
	}
	rsp, err := service.GetTcpClient().UserLogin(c, req)
	if err != nil {
		log.Error(fmt.Sprintf("User '%v' Login in failed,err:%v", userInfo.Username, err))
		c.JSON(constant.SystemError, gin.H{
			"msg": err,
		})
		return
	}

	if rsp.Status {
		log.Info(fmt.Sprintf("User '%v' login successfully", userInfo.Username))
		token := middleware.GenerateJwtToken(userInfo.Username)
		if token == "" {
			log.Info(fmt.Sprintf("User '%v' generate token failed", userInfo.Username))
			c.JSON(constant.SystemError, gin.H{
				"msg": "generate token failed",
			})
			return
		}
		cache.Rdb.Set(fmt.Sprintf("token_%v", userInfo.Username), token, 5*time.Minute)
		c.JSON(constant.Success, gin.H{
			"token": token,
		})
		return
	}
	log.Error(fmt.Sprintf("User '%v' login failed. err:%v", userInfo.Username, rsp.Msg))

	c.JSON(constant.InvalidUsernameOrPassword, gin.H{
		"msg": rsp.Msg,
	})
}
func Register(c *gin.Context) {
	log.Info("User Invoke Register API")
	defer func() {
		log.Info("User Finish  Register API")
	}()
	userInfo := UserParam{}
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		log.Error(fmt.Sprintf("Register failed,err:%v", err))

		c.JSON(constant.InvalidParam, gin.H{
			"msg": "Register failed,please check param",
		})
		return
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(userInfo.Password))
	encryptPassword := hex.EncodeToString(md5Ctx.Sum(nil))
	req := &api.RegisterUserReq{
		UserId:   userInfo.Username,
		Password: encryptPassword,
		NickName: userInfo.Nickname,
	}
	log.Info(fmt.Sprintf("User '%v' Invloke TcpClient:RegisterUser", userInfo.Username))
	rsp, err := service.GetTcpClient().RegisterUser(c, req)
	if err != nil {
		log.Error(fmt.Sprintf("User '%v' Register failed,err:%v", userInfo.Username, err))
		c.JSON(constant.SystemError, gin.H{
			"msg": fmt.Sprintf("register falid! err :%v", err),
		})
		return
	}
	if rsp.Status {
		log.Info(fmt.Sprintf("User '%v' Register successfully", userInfo.Username))
		c.JSON(constant.Success, gin.H{
			"msg": "register success! pleas login",
		})
		return
	}
	log.Error(fmt.Sprintf("User '%v' Register failed,err:%v", userInfo.Username, rsp.Msg))
	c.JSON(constant.UserExist, gin.H{
		"msg": rsp.Msg,
	})

}
func UpdateUserInfo(c *gin.Context) {
	log.Info("User Invoke UpdateUserInfo API")
	defer func() {
		log.Info("User Finish UpdateUserInfo API")
	}()
	userInfo := UpdateUserMask{}
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		log.Error(fmt.Sprintf("phase update json failed,err:%v", err))
		c.JSON(constant.InvalidParam, gin.H{
			"msg": "failed,please check param:",
		})
		return
	}
	name, exist := c.Get("username")
	if !exist {
		log.Info("system error! please login again")
		c.JSON(constant.SystemError, gin.H{
			"msg": "system error! please login again",
		})
		return
	}
	userInfo.Username = name.(string)
	oldEncryptPassword := ""
	if userInfo.OldPassword != "" {
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(userInfo.OldPassword))
		oldEncryptPassword = hex.EncodeToString(md5Ctx.Sum(nil))
	}
	newEncryptPassword := ""
	if userInfo.NewPassword != "" {
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(userInfo.NewPassword))
		newEncryptPassword = hex.EncodeToString(md5Ctx.Sum(nil))
	}

	req := &api.UserEditReq{
		UserId:      userInfo.Username,
		NewPassword: newEncryptPassword,
		OldPassword: oldEncryptPassword,
		NickName:    userInfo.NickName,
	}
	log.Info(fmt.Sprintf("User '%v' Invloke TcpClient:EditUserInfo", userInfo.Username))
	rsp, err := service.GetTcpClient().EditUserInfo(c, req)
	if err != nil {
		log.Error(fmt.Sprintf("User '%v' update failed,err:%v", userInfo.Username, err))
		c.JSON(constant.SystemError, gin.H{
			"msg": fmt.Sprintf("update falid! err :%v", err),
		})
		return
	}
	if rsp.Status {
		log.Info(fmt.Sprintf("User '%v' Update successfully,Del token cache", userInfo.Username))
		cache.Rdb.Del(fmt.Sprintf("token_%v", userInfo.Username)).Result()
		c.JSON(constant.Success, gin.H{
			"msg": "Edit UserInfo success! pleas login again",
		})
		return
	}
	log.Error(fmt.Sprintf("User '%v' Update failed,err:%v", userInfo.Username, rsp.Msg))
	c.JSON(constant.UpdateFailed, gin.H{
		"msg": rsp.Msg,
	})

}
func GetUserInfo(c *gin.Context) {
	log.Info("User Invoke GetUserInfo API")
	defer func() {
		log.Info("User Finish GetUserInfo API")
	}()
	userName, exist := c.Get("username")
	if !exist {
		log.Info("system error! please login again")
		c.JSON(constant.SystemError, gin.H{
			"msg": "system error! please login again",
		})
		return
	}
	req := &api.UserInfoReq{UserId: userName.(string)}
	log.Info(fmt.Sprintf("User '%v' Invloke TcpClient:UserInfo", userName.(string)))
	rsp, err := service.GetTcpClient().UserInfo(c, req)
	if err != nil {
		log.Error(fmt.Sprintf("User '%v' get userinfo failed,err:%v", userName.(string), err))
		c.JSON(constant.SystemError, gin.H{
			"msg": fmt.Sprintf("Get userInfo falid! err :%v", err),
		})
		return
	}
	log.Info(fmt.Sprintf("User '%v' Get info successfully", userName))
	profile := readPicture(rsp.ProfilePicture)
	c.JSON(constant.Success, gin.H{
		"username":        rsp.UserId,
		"nickName":        rsp.NickName,
		"profile_picture": profile,
	})

}
func UploadPicture(c *gin.Context) {
	log.Info("User Invoke UploadPicture API")
	defer func() {
		log.Info("User Finish UploadPicture API")
	}()

	name, exist := c.Get("username")
	if !exist {
		log.Info("system error! please login again")
		c.JSON(constant.SystemError, gin.H{
			"msg": "system error! please login again",
		})
		return
	}
	file, _ := c.FormFile("image")
	pictureUrl := path.Join("./avator", name.(string)) + ".jpg"
	err := c.SaveUploadedFile(file, pictureUrl)
	if err != nil {
		log.Error(fmt.Sprintf("user '%v' save picture failed,err : %v", name.(string), err))
		c.JSON(constant.InvalidData, gin.H{
			"msg": "upload picture failed",
		})
		return
	}
	req := &api.UploadProfileReq{
		UserId: name.(string),
		Url:    pictureUrl,
	}
	log.Info(fmt.Sprintf("User '%v' Invloke TcpClient:UploadProfile", name.(string)))
	rsp, err := service.GetTcpClient().UploadProfile(c, req)
	if err != nil {
		log.Error(fmt.Sprintf("User '%v' upload failed,err:%v", name.(string), err))
		c.JSON(constant.SystemError, gin.H{
			"msg": "upload failed",
		})
		return
	}
	if rsp.Status {
		log.Info(fmt.Sprintf("User '%v' upload  successfully", name.(string)))
		c.JSON(constant.Success, gin.H{
			"msg": "upload success",
		})
		return
	}
	log.Error(fmt.Sprintf("User '%v' upload  failed", name.(string)))
	c.JSON(constant.InvalidData, gin.H{
		"msg": rsp.Msg,
	})

}

// 身份认证
func Auth(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {
		c.JSON(constant.UnAuthorized, gin.H{
			"msg": "Please Login first",
		})
		c.Abort()
		return
	}

	name := middleware.PhaseJwtToken(token)
	if name == "" {
		c.JSON(constant.UnAuthorized, gin.H{
			"msg": "Invalid token",
		})
		c.Abort()
		return
	}

	c.Set("username", name)
	val, err := cache.Rdb.Get(fmt.Sprintf("token_%v", name)).Result()
	if err != nil && err != redis.Nil {
		log.Error(fmt.Sprintf("get cache error,err:%v", err))
		c.JSON(constant.UnAuthorized, gin.H{
			"msg": "System error",
		})
		c.Abort()
		return
	}

	if err == redis.Nil || val != token {
		// 缓存中没有或者 token不匹配
		c.JSON(constant.UnAuthorized, gin.H{
			"msg": "Please Login again",
		})
		c.Abort()
		return
	}
	cache.Rdb.Expire(fmt.Sprintf("token_%v", name), 5*time.Minute)
	c.Next()
}

// 图片解码存储
func DecodePic(base64String string, username string) string {
	//ff, _ := ioutil.ReadFile("./下载.jpeg")
	//base64Data := base64.StdEncoding.EncodeToString(ff)
	//fmt.Println(base64Data)
	////解压
	dist, _ := base64.StdEncoding.DecodeString(base64String)
	//写入新文件
	_ = ioutil.WriteFile(fmt.Sprintf("./%v.jpg", username), dist, os.ModePerm)
	return ""
}

func readPicture(url string) string {
	_, err := os.Stat(url)
	if err != nil {
		// no such file or dir
		log.Error(fmt.Sprintf("No such File:%v", url))
		return ""
	}
	ff, _ := ioutil.ReadFile(url)
	base64Data := base64.StdEncoding.EncodeToString(ff)
	return base64Data
}
