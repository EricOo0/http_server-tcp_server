package router

import (
	"github.com/gin-gonic/gin"
	"http_server/src/handler"
	"http_server/src/service"
	"net/http"
)

func InitRouter() {
	engine := service.GetHttpServer()
	//模板解析
	engine.LoadHTMLGlob("template/*")
	// 首页
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	// 注册页面
	engine.GET("/userRegister", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{})
	})

	// 用户信息
	engine.GET("/userinfo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "userinfo.tmpl", gin.H{})
	})

	engine.GET("/updateUser", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update.tmpl", gin.H{})
	})

	// 登陆
	engine.POST("/login", handler.Login)
	// 注册
	engine.POST("/register", handler.Register)
	//修改
	engine.POST("/update", handler.Auth, handler.UpdateUserInfo)
	//获取profile
	engine.GET("/profile", handler.Auth, handler.GetUserInfo)
	//上传头像
	engine.POST("/upload", handler.Auth, handler.UploadPicture)

}
