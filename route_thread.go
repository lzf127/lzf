package main

import (
	"fmt"
	"net/http"
	"time"

	"./data"
)

//newThread 新建帖子载入
func newThread(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		shtml(w, "login.layout", "private.navbar", "new.thread")
	}
}

//createThread 新建帖子
func createThread(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		th := r.PostFormValue("topic")
		user, err := sess.User()
		if err != nil {
			fmt.Println("获取用户失败")
		}
		_, err = user.AddThreads(th)
		if err != nil {
			fmt.Println("创建帖子失败11")
		}
		http.Redirect(w, r, "/", 302)
	}
}

// type Post struct {
// 	Id        int
// 	Uuid      string
// 	Body      string
// 	UserId    int
// 	ThreadId  int
// 	CreatedAt time.Time
// }
func postThread(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			fmt.Println("sess查找失败")
		}
		p := data.Post{}
		p.Body = r.PostFormValue("body")
		p.Uuid = r.PostFormValue("uuid")
		th, _ := data.TUuid(p.Uuid)
		p.ThreadId = th.Id
		p.UserId = user.Id
		p.CreatedAt = time.Now()
		err = data.Cpost(p)
		if err != nil {
			fmt.Println("创建评价失败")
		}
		url := fmt.Sprint("/thread/read?id=", p.Uuid)
		http.Redirect(w, r, url, 302)
	}

}

func readThread(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	uuid := vals.Get("id")
	t, err := data.TUuid(uuid)
	if err != nil {
		fmt.Println("获取帖子失败")
	}
	_, err = session(w, r)
	if err != nil {
		hhtml(w, &t, "layout", "public.navbar", "public.thread")
	} else {
		hhtml(w, &t, "layout", "private.navbar", "private.thread")
	}
}
