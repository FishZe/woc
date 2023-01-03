package main

import (
	"fmt"
	"log"
	"testing"
)

func TestInitDB(t *testing.T) {
	err := InitDB()
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsertUser(t *testing.T) {
	TestInitDB(t)
	err := InsertUser(USER{UserName: "admin", Password: "test", Email: "123@123.com", Role: 1})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSomeUsers(t *testing.T) {
	TestInitDB(t)
	users := GetSomeUsers(0, 10)
	if len(users) == 0 {
		t.Fatal("GetSomeUsers failed")
	}
	log.Printf("%v", users)
}

func TestSearchUser(t *testing.T) {
	TestInitDB(t)
	users := SearchUser(USER{UserName: "陈睿", Role: -2})
	if len(users) == 0 {
		t.Fatal("SearchUser failed")
	}
	log.Printf("%v", users)
}

func TestDeleteUser(t *testing.T) {
	TestInitDB(t)
	users := DeleteUser(USER{UserName: "陈睿", Role: -2})
	fmt.Println(users)
}

func TestModifyUserById(t *testing.T) {
	TestInitDB(t)
	err := ModifyUserById(USER{Id: 1, Email: "1234@123.com"})
	if err != nil {
		t.Fatal(err)
	}
}
