package data

import (
	"errors"
	"fmt"
	"time"
)

//User 用户
type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

//Session 会话结构
type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

//Check 检查数据库中是否存在
func (se *Session) Check() (va bool, err error) {
	sts := "select id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1"
	err = Db.QueryRow(sts, se.Uuid).Scan(&se.Id, &se.Uuid, &se.Email, &se.UserId, &se.CreatedAt)
	if err != nil {
		va = false
		return
	}
	if se.Id != 0 {
		va = true
	}
	return
}

//SearchEmail 通过邮箱查找用户
func SearchEmail(email string) (user User, err error) {
	ss, err := Db.Query("select id, uuid, name, email, password, created_at from users where email = $1", email)
	if err != nil {
		fmt.Println("查找用户错误")
		return
	}
	ss.Next()
	defer ss.Close()
	err = ss.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

//SearchUuid 通过Uuid查找用户
func SearchUuid(uuid string) (user User, err error) {
	ss, err := Db.Query("select id, uuid, name, email, password, created_at from users where uuid = $1", uuid)
	if err != nil {
		fmt.Println("查找用户错误")
		return
	}
	ss.Next()
	defer ss.Close()
	err = ss.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

//CreateSession 创建会话
func (user *User) CreateSession() (ses Session, err error) {
	sta := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
	stmt, err := Db.Prepare(sta)
	if err != nil {
		fmt.Println("创建会话失败")
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(createUUID(), user.Email, user.Id, time.Now()).Scan(&ses.Id, &ses.Uuid, &ses.Email, &ses.UserId, &ses.CreatedAt)
	return
}

//CreateUser 添加用户
func CreateUser(user *User) error {
	ss, err := SearchEmail(user.Email)
	if ss.Email != user.Email {
		sta := "insert into users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5)"
		stmt, err := Db.Prepare(sta)
		if err != nil {
			fmt.Println("新建用户失败")
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(createUUID(), user.Name, user.Email, user.Password, time.Now())
		if err != nil {
			fmt.Println("新建用户添加数据库失败")
			return err
		}
		return err
	}
	fmt.Println("用户已经存在")
	err = errors.New("用户已经存在")
	return err
}
