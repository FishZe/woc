package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserInfo(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("USER_ID") == nil {
		returnMsg(c, false, http.StatusUnauthorized, UserNotLoginMsg, nil)
	} else {
		nowUser := SearchUser(USER{Id: session.Get("USER_ID").(int), Role: -2})
		if len(nowUser) == 0 {
			returnMsg(c, false, http.StatusUnauthorized, UserNotExistMsg, nil)
		} else {
			returnMsg(c, true, http.StatusOK, NoMsg, nowUser[0])
		}
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("USER_ID") == nil {
		returnMsg(c, false, http.StatusUnauthorized, UserNotLoginMsg, nil)
	} else {
		session.Delete("USER_ID")
		session.Delete("USER_ROLE")
		err := session.Save()
		if err != nil {
			returnMsg(c, false, http.StatusInternalServerError, SessionErrorMsg, nil)
		} else {
			returnMsg(c, true, http.StatusOK, NoMsg, nil)
		}
	}
}

// Login is a function for user login
func Login(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil || u.UserName == "" || u.Password == "" {
		returnMsg(c, false, http.StatusBadRequest, BadRequestMsg, nil)
		return
	}
	if LoginUser(u) {
		nowUser := SearchUser(USER{UserName: u.UserName, Role: -2})
		if nowUser[0].Role != 1 {
			returnMsg(c, false, http.StatusUnauthorized, UserNotAdminMsg, nil)
			return
		}
		session := sessions.Default(c)
		session.Set("USER_ID", nowUser[0].Id)
		session.Set("USER_ROLE", nowUser[0].Role)
		err = session.Save()
		if err != nil {
			returnMsg(c, false, http.StatusInternalServerError, SessionErrorMsg, nil)
		} else {
			returnMsg(c, true, http.StatusOK, NoMsg, nil)
		}
	} else {
		returnMsg(c, false, http.StatusUnauthorized, UserNotExistMsg, nil)
	}
}

// AdminNewUser is a function for admin create a new user
func AdminNewUser(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil || u.UserName == "" || u.Password == "" {
		returnMsg(c, false, http.StatusBadRequest, BadRequestMsg, nil)
		return
	}
	// Check user exist
	nowUser := SearchUser(USER{UserName: u.UserName, Role: -2})
	if len(nowUser) != 0 {
		returnMsg(c, false, http.StatusBadRequest, UserNameExistMsg, nil)
	} else {
		err = InsertUser(u)
		if err != nil {
			returnMsg(c, false, http.StatusInternalServerError, InternalServerErrorMsg, nil)
		} else {
			returnMsg(c, true, http.StatusOK, NoMsg, nil)
		}
	}
}

// AdminDeleteUser is a function for admin delete a user
func AdminDeleteUser(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil || u.Id <= 0 {
		returnMsg(c, false, http.StatusBadRequest, BadRequestMsg, nil)
		return
	}
	nowUser := SearchUser(USER{Id: u.Id, Role: -2})
	if len(nowUser) == 0 {
		returnMsg(c, false, http.StatusBadRequest, UserNotExistMsg, nil)
		return
	} else {
		session := sessions.Default(c)
		if session.Get("USER_ID").(int) == u.Id {
			returnMsg(c, false, http.StatusInternalServerError, "请勿花样作死", nil)
			return
		}
		err = DeleteUser(USER{Id: u.Id})
		if err != nil {
			returnMsg(c, false, http.StatusInternalServerError, InternalServerErrorMsg, nil)
		} else {
			returnMsg(c, true, http.StatusOK, NoMsg, nil)
		}
	}
}

// AdminSearchUser is a function for admin update a user
func AdminSearchUser(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil {
		returnMsg(c, false, http.StatusBadRequest, BadRequestMsg, nil)
	} else {
		users := SearchUser(u)
		var data []USER
		for _, v := range users {
			data = append(data, v)
		}
		returnMsg(c, true, http.StatusOK, NoMsg, data)
	}
}

// AdminUpdateUser is a function for admin update a user
func AdminUpdateUser(c *gin.Context) {
	u := USER{}
	err := c.BindJSON(&u)
	if err != nil || u.Id <= 0 {
		returnMsg(c, false, http.StatusBadRequest, BadRequestMsg, nil)
		return
	}
	nowUser := SearchUser(USER{Id: u.Id, Role: -2})
	if len(nowUser) == 0 {
		returnMsg(c, false, http.StatusBadRequest, UserNotExistMsg, nil)
		return
	}
	nowUser = SearchUser(USER{UserName: u.UserName, Role: -2})
	if len(nowUser) != 0 && nowUser[0].Id != u.Id {
		returnMsg(c, false, http.StatusBadRequest, UserNameExistMsg, nil)
		return
	}
	user := USER{Id: u.Id, UserName: u.UserName, Password: u.Password, Email: u.Email, Role: u.Role}
	err = ModifyUserById(user)
	if err != nil {
		returnMsg(c, false, http.StatusInternalServerError, InternalServerErrorMsg, nil)
	} else {
		returnMsg(c, true, http.StatusOK, NoMsg, nil)
	}
}

// AdminGetSomeUser is a function for admin get all user
func AdminGetSomeUser(c *gin.Context) {
	type UserGetSome struct {
		FromId int `json:"from"`
		Sum    int `json:"sum" binding:"required"`
	}
	u := UserGetSome{}
	err := c.BindJSON(&u)
	if err != nil || u.FromId < 0 || u.Sum <= 0 {
		returnMsg(c, false, http.StatusBadRequest, BadRequestMsg, nil)
	} else {
		users, allSum := GetSomeUsers(u.FromId, u.Sum)
		var data []USER
		for _, v := range users {
			data = append(data, v)
		}
		res := make(map[string]interface{})
		res["user"] = data
		res["all"] = allSum
		page := allSum / int64(u.Sum)
		if allSum%int64(u.Sum) != 0 {
			page++
		}
		res["page"] = page
		returnMsg(c, true, http.StatusOK, "", res)
	}
}
