package data

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	//数据库包
	_ "github.com/lib/pq"
)

//Db 数据库连接
var Db *sql.DB

//Post 评价
type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

//连接数据库
func init() {
	var err error
	Db, err = sql.Open("postgres", "user=lzf dbname=chitchat password=123456 sslmode=disable")
	if err != nil {
		fmt.Println("数据库连接错误")
	}
	return
}

//createUUID 创建UUID
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		fmt.Println("创建UUID失败")
	}
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}
