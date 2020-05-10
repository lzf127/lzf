package main

import (
	"fmt"
	"net/http"

	"./data"
)

//login 载入login界面
func login(w http.ResponseWriter, r *http.Request) {
	shtml(w, "login.layout", "public.navbar", "login")
}

//logout 退出登入
func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		ses := data.Session{Uuid: cookie.Value}
		ses.Delete()
	}
	http.Redirect(w, r, "/", 302)
}

//signup 载入signup界面
func signup(w http.ResponseWriter, r *http.Request) {
	shtml(w, "login.layout", "public.navbar", "signup")
}

//signupAccount 创建用户
func signupAccount(w http.ResponseWriter, r *http.Request) {
	user := data.User{}
	user.Name = r.FormValue("name")
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	err := data.CreateUser(&user)
	if err != nil {
		fmt.Println("用户创建失败")
		http.Redirect(w, r, "/signup", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}

}

//authenticate 用户登入验证并加载cookie
func authenticate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := data.SearchEmail(email)
	if err != nil {
		fmt.Println("用户不存在")
		http.Redirect(w, r, "/login", 302)
	} else {
		if user.Password != password {
			fmt.Println("密码错误")
			http.Redirect(w, r, "/login", 302)
		} else {
			ss, err := user.CreateSession()
			if err != nil {
				fmt.Println("创建会话失败2")
			}
			c := http.Cookie{
				Name:     "_cookie",
				Value:    ss.Uuid,
				HttpOnly: true,
			}
			http.SetCookie(w, &c)
			http.Redirect(w, r, "/", 302)
		}
	}
}
