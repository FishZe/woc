package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	router := gin.Default()
	err := InitDB()
	if err != nil {
		return
	}
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("SESSION", store))
	router.POST("/login", Login)
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
		log.Fatalf("Gin start error: %v", err)
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("USER_ID") == nil {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "data": map[string]string{"msg": "not login"}})
			c.Abort()
		} else if session.Get("USER_ROLE") != 1 {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "data": map[string]string{"msg": "not admin"}})
			c.Abort()
		}
		c.Next()
	}
}

func Login(c *gin.Context) {
	type UserLogin struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	u := UserLogin{}
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": "get user info error:" + err.Error()}})
		return
	}
	if u.Name == "" || u.Password == "" {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": map[string]string{"msg": "name or password is empty"}})
		return
	}
	user := USER{UserName: u.Name, Password: u.Password}
	if LoginUser(user) {
		nowUser := SearchUser(USER{UserName: u.Name, Role: -2})
		session := sessions.Default(c)
		session.Set("USER_ID", nowUser[0].Id)
		session.Set("USER_ROLE", nowUser[0].Role)
		err = session.Save()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": "session save error"}})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": map[string]string{"msg": "login success"}})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "data": map[string]string{"msg": "name or password is wrong"}})
	}
}

func AdminNewUser(c *gin.Context) {
	type UserNew struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Role     int    `json:"role"`
	}
	u := UserNew{}
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": "get new user info error: " + err.Error()}})
		return
	}
	if u.Name == "" || u.Password == "" {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": map[string]string{"msg": "name or password is empty"}})
		return
	}
	nowUser := SearchUser(USER{UserName: u.Name, Role: -2})
	if len(nowUser) != 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": map[string]string{"msg": "user already exist"}})
		return
	}
	user := USER{UserName: u.Name, Password: u.Password, Email: u.Email, Role: u.Role}
	err = InsertUser(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": "insert user into database error"}})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": map[string]string{"msg": "insert user success"}})
	}
}

func AdminDeleteUser(c *gin.Context) {
	type UserDelete struct {
		Id int `json:"id" binding:"required"`
	}
	u := UserDelete{}
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": "get delete user info error: " + err.Error()}})
		return
	}
	if u.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": map[string]string{"msg": "id is empty"}})
		return
	}
	nowUser := SearchUser(USER{Id: u.Id, Role: -2})
	if len(nowUser) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": map[string]string{"msg": "user not exist"}})
		return
	}
	err = DeleteUser(USER{Id: u.Id})
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": map[string]string{"msg": "delete user success"}})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": err.Error()}})
	}
}

func AdminSearchUser(c *gin.Context) {
	type UserSearch struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     int    `json:"role"`
	}
	u := UserSearch{}
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": "get search user info error:" + err.Error()}})
		return
	}
	users := SearchUser(USER{Id: u.Id, UserName: u.Name, Password: u.Password, Email: u.Email, Role: u.Role})
	var data []UserSearch
	for _, v := range users {
		data = append(data, UserSearch{Id: v.Id, Name: v.UserName, Password: v.Password, Email: v.Email, Role: v.Role})
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})
}

func AdminUpdateUser(c *gin.Context) {
	type UserUpdate struct {
		Id       int    `json:"id" binding:"required"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     int    `json:"role"`
	}
	u := UserUpdate{}
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": "get update user info error:" + err.Error()}})
		return
	}
	if u.Id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": map[string]string{"msg": "id is empty"}})
		return
	}
	nowUser := SearchUser(USER{Id: u.Id, Role: -2})
	if len(nowUser) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": map[string]string{"msg": "user not exist"}})
		return
	}
	nowUser = SearchUser(USER{UserName: u.Name, Role: -2})
	if len(nowUser) != 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": map[string]string{"msg": "user name already exist"}})
		return
	}
	user := USER{Id: u.Id, UserName: u.Name, Password: u.Password, Email: u.Email, Role: u.Role}
	err = ModifyUserById(user)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": map[string]string{"msg": "update user success"}})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": err.Error()}})
	}
}

func AdminGetSomeUser(c *gin.Context) {
	type UserGetSome struct {
		FromId int `json:"from_id"`
		Sum    int `json:"sum"`
	}
	u := UserGetSome{}
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": map[string]string{"msg": "get json info error:" + err.Error()}})
		return
	}
	users := GetSomeUsers(u.FromId, u.Sum)
	type User struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     int    `json:"role"`
	}
	var data []User
	for _, v := range users {
		data = append(data, User{Id: v.Id, Name: v.UserName, Password: v.Password, Email: v.Email, Role: v.Role})
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})
}
