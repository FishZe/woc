package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type USER struct {
	Id       int    `json:"id"`
	UserName string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
}

// InitDB init database
func InitDB() error {
	err := error(nil)
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&USER{})
	if err != nil {
		return err
	}
	return nil
}

// InsertUser insert a user
func InsertUser(user USER) error {
	err := DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser delete a user
func DeleteUser(user USER) error {
	err := DB.Where("id = ?", user.Id).Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}

// LoginUser login user
func LoginUser(user USER) bool {
	var users []USER
	DB.Where("user_name = ? AND password = ?", user.UserName, user.Password).Find(&users)
	if len(users) == 0 {
		return false
	}
	return true
}

// GetSomeUsers get some users
func GetSomeUsers(start, sum int) []USER {
	var users []USER
	DB.Limit(sum).Offset(start).Find(&users)
	return users
}

// SearchUser search user
func SearchUser(user USER) []USER {
	var users []USER
	db := DB
	if user.Id != 0 {
		db = db.Where("id = ?", user.Id).Find(&users)
	}
	if user.UserName != "" {
		db = db.Where("user_name = ?", user.UserName).Find(&users)
	}
	if user.Password != "" {
		db = db.Where("password = ?", user.Password).Find(&users)
	}
	if user.Email != "" {
		db = db.Where("email = ?", user.Email).Find(&users)
	}
	if user.Role != -2 {
		db = db.Where("role = ?", user.Role).Find(&users)
	}
	return users
}

// ModifyUserById modify user
func ModifyUserById(user USER) error {
	err := DB.Where("id = ?", user.Id).Updates(user).Error
	if err != nil {
		return err
	}
	return nil
}
