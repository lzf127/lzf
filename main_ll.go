package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"./data"
)

//index 界面并检查是否登入
func index(w http.ResponseWriter, r *http.Request) {
	threads, err := data.Threads()
	if err != nil {
		fmt.Println("获取帖子失败")
	} else {
		_, err = session(w, r)
		if err != nil {
			hhtml(w, threads, "layout", "public.navbar", "index")
		} else {
			hhtml(w, threads, "layout", "private.navbar", "index")
		}
	}
}

//合并html并添加帖子
func hhtml(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

//session检查用户是否登入
func session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = data.Session{}
		sess.Uuid = cookie.Value
		ok, _ := sess.Check()
		if !ok {
			//fmt.Println(ok)
			err = errors.New("没有 session")
		}
	}
	//fmt.Println(cookie.Value)
	return
}

//合并hmtl
func shtml(w http.ResponseWriter, filenames ...string) {
	var files []string
	t := template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	err := t.Execute(w, nil)
	if err != nil {
		fmt.Println("输出错误")
	}
}
