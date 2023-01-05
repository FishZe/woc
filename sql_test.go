package main

import (
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"math/rand"
	"testing"
	"time"
)

var randUsers []USER

func randStr(length int) string {
	bytes := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%^&*()_+{},./;'[]<>?:")
	var result []byte
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func mkRandUser() USER {
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	return USER{
		UserName: randStr(rand.Intn(10) + 1),
		Password: randStr(rand.Intn(10) + 1),
		Email:    randStr(rand.Intn(10)+1) + "@" + randStr(rand.Intn(5)+1) + "." + randStr(rand.Intn(5)+1),
		Role:     rand.Intn(3) - 1,
	}
}

func MkRandUsers() {
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	sum := rand.Intn(1000) + 1000
	userName := make(map[string]struct{})
	userEmail := make(map[string]struct{})
	for i := 0; i < sum; i++ {
		nowUser := mkRandUser()
		if _, ok := userName[nowUser.UserName]; !ok {
			if _, ok = userEmail[nowUser.Email]; !ok {
				randUsers = append(randUsers, nowUser)
				userName[nowUser.UserName] = struct{}{}
				userEmail[nowUser.Email] = struct{}{}
			}
		}

	}
	log.Printf("MkRandUsers success, sum: %d", sum)
}

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[Woc][%time%][%lvl%]: %msg% \n",
	})
}

func TestInitDB(t *testing.T) {
	err := InitDB()
	if err != nil {
		t.Fatal(err)
	}
	DropTable(t)
	log.Printf("Init DB success")
}

func DropTable(t *testing.T) {
	if DB == nil {
		TestInitDB(t)
	}
	err := DB.Where("1 = 1").Delete(&USER{}).Error
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Drop the table success")
}

func TestInsertUser(t *testing.T) {
	if DB == nil {
		TestInitDB(t)
	}
	MkRandUsers()
	for _, user := range randUsers {
		err := InsertUser(user)
		if err != nil {
			t.Fatal(err)
		}
	}
	log.Printf("InsertUser success, sum: %d", len(randUsers))
}

func TestGetSomeUsers(t *testing.T) {
	if DB == nil {
		TestInitDB(t)
	}
	for i := 0; i < rand.Intn(len(randUsers))+1; i++ {
		from := rand.Intn(len(randUsers))
		sum := rand.Intn(len(randUsers) - from)
		for sum == 0 {
			sum = rand.Intn(len(randUsers) - from)
		}
		users := GetSomeUsers(from, sum)
		if len(users) != sum {
			t.Fatal("GetSomeUsers error")
		}
	}
	log.Printf("GetSomeUsers success")
}

func TestSearchUser(t *testing.T) {
	if DB == nil {
		TestInitDB(t)
	}
	// test Id
	for i := 1; i < rand.Intn(len(randUsers))+1; i++ {
		user := SearchUser(USER{Id: i, Role: -2})
		if len(user) != 1 {
			t.Fatal("SearchUser error: ", randUsers[i-1])
		}
	}
	// test name
	for i := 0; i < rand.Intn(len(randUsers))+1; i++ {
		name := randUsers[rand.Intn(len(randUsers))].UserName
		users := SearchUser(USER{UserName: name, Role: -2})
		if len(users) != 1 {
			t.Fatal("SearchUser error: ", name)
		}
	}
	// test email
	for i := 0; i < rand.Intn(len(randUsers))+1; i++ {
		email := randUsers[rand.Intn(len(randUsers))].Email
		users := SearchUser(USER{Email: email, Role: -2})
		if len(users) != 1 {
			t.Fatal("SearchUser error: ", email)
		}
	}
	// test role
	for i := -1; i < 2; i++ {
		userSum := 0
		for _, user := range randUsers {
			if user.Role == i {
				userSum++
			}
		}
		users := SearchUser(USER{Role: i})
		if len(users) != userSum {
			t.Fatal("SearchUser error: ", i)
		}
	}
}

func TestLoginUser(t *testing.T) {
	if DB == nil {
		TestInitDB(t)
	}
	for _, v := range randUsers {
		if !LoginUser(v) {
			t.Fatal("LoginUser error: ", v)
		}
	}
	for i := 0; i < rand.Intn(len(randUsers))+1; i++ {
		user := mkRandUser()
		if LoginUser(user) {
			find := false
			for _, v := range randUsers {
				if v.UserName == user.UserName && v.Password == user.Password {
					find = true
					break
				}
			}
			if !find {
				t.Fatal("LoginUser error: ", user)
			}
		}
	}
	log.Printf("LoginUser success")
}

func TestModifyUserById(t *testing.T) {
	if DB == nil {
		TestInitDB(t)
	}
	for i := 0; i < rand.Intn(len(randUsers))+1; i++ {
		user := mkRandUser()
		user.Id = rand.Intn(len(randUsers))
		err := ModifyUserById(user)
		if err != nil {
			t.Fatal(err)
		}
	}
	log.Printf("ModifyUserById success")
}

func TestDeleteUser(t *testing.T) {
	if DB == nil {
		TestInitDB(t)
	}
	for i, v := range randUsers {
		v.Id = i
		err := DeleteUser(v)
		if err != nil {
			t.Fatal(err)
		}
	}
	log.Printf("DeleteUser success")
}
