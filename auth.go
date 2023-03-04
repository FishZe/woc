package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Auth is a middleware for check user admin
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("USER_ID") == nil {
			returnMsg(c, false, http.StatusUnauthorized, UserNotLoginMsg, nil)
		} else if session.Get("USER_ROLE") != 1 {
			returnMsg(c, false, http.StatusUnauthorized, UserNotAdminMsg, nil)
		}
		c.Next()
	}
}
