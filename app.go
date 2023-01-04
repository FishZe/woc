package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	NoMsg                  = ""
	UserNotExistMsg        = "user not exist"
	UserNotLoginMsg        = "User not login"
	UserNotAdminMsg        = "User not admin"
	UserNameExistMsg       = "user name already exist"
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
	// Admin
	admin := router.Group("/admin")
	{
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

func returnMsg(c *gin.Context, ReturnCode int, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": ReturnCode, "msg": msg, "data": data})
}

// Auth is a middleware for check user admin
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("USER_ID") == nil {
			returnMsg(c, http.StatusUnauthorized, UserNotLoginMsg, nil)
		} else if session.Get("USER_ROLE") != 1 {
			returnMsg(c, http.StatusUnauthorized, UserNotAdminMsg, nil)
		}
		c.Next()
	}
}

// Login is a function for user login
func Login(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil || u.UserName == "" || u.Password == "" {
		returnMsg(c, http.StatusBadRequest, BadRequestMsg, nil)
		return
	}
	if LoginUser(u) {
		nowUser := SearchUser(USER{UserName: u.UserName, Role: -2})
		session := sessions.Default(c)
		session.Set("USER_ID", nowUser[0].Id)
		session.Set("USER_ROLE", nowUser[0].Role)
		err = session.Save()
		if err != nil {
			returnMsg(c, http.StatusInternalServerError, SessionErrorMsg, nil)
		} else {
			returnMsg(c, http.StatusOK, NoMsg, nil)
		}
	} else {
		returnMsg(c, http.StatusUnauthorized, UserNotExistMsg, nil)
	}
}

// AdminNewUser is a function for admin create a new user
func AdminNewUser(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil || u.UserName == "" || u.Password == "" {
		returnMsg(c, http.StatusBadRequest, BadRequestMsg, nil)
		return
	}
	// Check user exist
	nowUser := SearchUser(USER{UserName: u.UserName, Role: -2})
	if len(nowUser) != 0 {
		returnMsg(c, http.StatusBadRequest, UserNameExistMsg, nil)
	} else {
		err = InsertUser(u)
		if err != nil {
			returnMsg(c, http.StatusInternalServerError, InternalServerErrorMsg, nil)
		} else {
			returnMsg(c, http.StatusOK, NoMsg, nil)
		}
	}
}

// AdminDeleteUser is a function for admin delete a user
func AdminDeleteUser(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil || u.Id <= 0 {
		returnMsg(c, http.StatusBadRequest, BadRequestMsg, nil)
		return
	}
	nowUser := SearchUser(USER{Id: u.Id, Role: -2})
	if len(nowUser) == 0 {
		returnMsg(c, http.StatusBadRequest, UserNotExistMsg, nil)
		return
	} else {
		err = DeleteUser(USER{Id: u.Id})
		if err != nil {
			returnMsg(c, http.StatusInternalServerError, InternalServerErrorMsg, nil)
		} else {
			returnMsg(c, http.StatusOK, NoMsg, nil)
		}
	}
}

// AdminSearchUser is a function for admin update a user
func AdminSearchUser(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil {
		returnMsg(c, http.StatusBadRequest, BadRequestMsg, nil)
	} else {
		users := SearchUser(u)
		var data []USER
		for _, v := range users {
			data = append(data, v)
		}
		returnMsg(c, http.StatusOK, "", data)
	}
}

// AdminUpdateUser is a function for admin update a user
func AdminUpdateUser(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil || u.Id <= 0 {
		returnMsg(c, http.StatusBadRequest, BadRequestMsg, nil)
		return
	}
	nowUser := SearchUser(USER{Id: u.Id, Role: -2})
	if len(nowUser) == 0 {
		returnMsg(c, http.StatusBadRequest, UserNotExistMsg, nil)
		return
	}
	nowUser = SearchUser(USER{UserName: u.UserName, Role: -2})
	if len(nowUser) != 0 {
		returnMsg(c, http.StatusBadRequest, UserNameExistMsg, nil)
		return
	}
	user := USER{Id: u.Id, UserName: u.UserName, Password: u.Password, Email: u.Email, Role: u.Role}
	err = ModifyUserById(user)
	if err != nil {
		returnMsg(c, http.StatusInternalServerError, InternalServerErrorMsg, nil)
	} else {
		returnMsg(c, http.StatusOK, NoMsg, nil)
	}
}

// AdminGetSomeUser is a function for admin get all user
func AdminGetSomeUser(c *gin.Context) {
	type UserGetSome struct {
		FromId int `json:"from_id" binding:"required"`
		Sum    int `json:"sum" binding:"required"`
	}
	u := UserGetSome{}
	err := c.BindJSON(&u)
	if err != nil || u.FromId <= 0 || u.Sum <= 0 {
		returnMsg(c, http.StatusBadRequest, BadRequestMsg, nil)
	} else {
		users := GetSomeUsers(u.FromId, u.Sum)
		var data []USER
		for _, v := range users {
			data = append(data, v)
		}
		returnMsg(c, http.StatusOK, "", data)
	}
}
