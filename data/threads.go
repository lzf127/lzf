package data

import (
	"fmt"
	"time"
)

//Thread ...贴子类型
type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

//Threads ...获取数据库帖子
func Threads() (t []Thread, err error) {
	rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")
	if err != nil {
		fmt.Println("数据库获取帖子失败")
	}
	for rows.Next() {
		tx := Thread{}
		if err = rows.Scan(&tx.Id, &tx.Uuid, &tx.Topic, &tx.UserId, &tx.CreatedAt); err != nil {
			return
		}
		t = append(t, tx)
	}
	defer rows.Close()
	return t, err
}

//AddThreads 创建帖子
func (user *User) AddThreads(s string) (th Thread, err error) {
	th.Topic = s
	th.UserId = user.Id
	th.Uuid = createUUID()
	th.CreatedAt = time.Now()
	sts := "insert into threads (uuid, topic, user_id, created_at) values ($1, $2, $3, $4) returning id"
	stmt, err := Db.Prepare(sts)
	if err != nil {
		fmt.Println("创建表数据库连接失败")
		return
	}
	defer stmt.Close()
	stmt.QueryRow(&th.Uuid, &th.Topic, &th.UserId, &th.CreatedAt).Scan(&th.Id)
	return
}

//User 通过会话查找用户userd的email
func (ses Session) User() (user User, err error) {
	err = Db.QueryRow("select email FROM sessions WHERE uuid = $1", ses.Uuid).
		Scan(&ses.Email)
	if err != nil {
		fmt.Println("数据库查找session错误")
	}
	user, err = SearchEmail(ses.Email)
	if err != nil {
		fmt.Println("用户错误")
		return
	}
	return
}

//Delete 表session数据
func (ses *Session) Delete() error {
	st := "delete from sessions where uuid =$1 "
	stm, err := Db.Prepare(st)
	if err != nil {
		fmt.Println("删除sessions表数据错误")
		return err
	}
	_, err = stm.Exec(ses.Uuid)
	if err != nil {
		fmt.Println("删除sessions表数据错误2")
	}
	return err
}

//CreatedAtDate 通过帖子方法获得时间
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

//NumReplies 通过帖子方法获得评价数
func (thread *Thread) NumReplies() (count int) {
	rows, err := Db.Query("SELECT count(*) FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		fmt.Println("获取评价数失败")
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	return
}

//User 通过帖子获得用户
func (thread *Thread) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", thread.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

//TUuid 通过uuid获得帖子
func TUuid(uuid string) (t Thread, err error) {
	th, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads where uuid=$1", uuid)
	if err != nil {
		fmt.Println("uuid访问数据库失败")
		return
	}
	th.Next()
	defer th.Close()
	err = th.Scan(&t.Id, &t.Uuid, &t.Topic, &t.UserId, &t.CreatedAt)
	if err != nil {
		fmt.Println("uid访问数据库失败2")
	}
	return
}

//Posts 通过帖子获得评价
func (t *Thread) Posts() (p []Post, err error) {
	stm := "SELECT id, uuid, body, user_id, thread_id, created_at FROM posts where thread_id=$1"
	th, err := Db.Query(stm, t.Id)
	if err != nil {
		fmt.Println("posts访问数据库失败")
		return
	}
	for th.Next() {
		pp := Post{}
		err = th.Scan(&pp.Id, &pp.Uuid, &pp.Body, &pp.UserId, &pp.ThreadId, &pp.CreatedAt)
		p = append(p, pp)
	}
	th.Close()
	return
}

//User 通过Post获取用户
func (p *Post) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", p.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

//(p *Post) CreatedAtDate() 通过post获得
func (p *Post) CreatedAtDate() string {
	return p.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

//Cpost 新增评价
func Cpost(p Post) error {
	sts := "insert into posts (uuid, body, user_id, thread_id, created_at) values ($1, $2, $3, $4, $5) returning id"
	stmt, err := Db.Prepare(sts)
	if err != nil {
		fmt.Println("创建表数据库连接失败")
		return err
	}
	defer stmt.Close()
	stmt.QueryRow(&p.Uuid, &p.Body, &p.UserId, &p.ThreadId, &p.CreatedAt).Scan(&p.Id)
	return err
}
