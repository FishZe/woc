package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	NoMsg                  = ""
	UserNotExistMsg        = "User not exist"
	UserNotLoginMsg        = "User not login"
	UserNotAdminMsg        = "User not admin"
	UserNameExistMsg       = "User name already exist"
	BadRequestMsg          = "Bad request"
	SessionErrorMsg        = "Session error"
	InternalServerErrorMsg = "Internal server error"
)

func Run() error {
	// Init DB
	err := InitDB()
	if err != nil {
		return err
	}
	// Route
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Session
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("SESSION", store))
	router.POST("/login", Login)
	router.POST("/logout", Logout)
	router.POST("/me", UserInfo)
	// Admin
	admin := router.Group("/admin")
	{
		// 需要鉴权的接口
		admin.POST("/new", Auth(), AdminNewUser)
		admin.POST("/delete", Auth(), AdminDeleteUser)
		admin.POST("/update", Auth(), AdminUpdateUser)
		admin.POST("/search", Auth(), AdminSearchUser)
		admin.POST("/get", Auth(), AdminGetSomeUser)
	}
	err = router.Run()
	if err != nil {
		return err
	}
	return nil
}

func returnMsg(c *gin.Context, success bool, errCode int, errMsg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"errCode": errCode,
		"errMsg":  errMsg,
		"data":    data,
	})
}
