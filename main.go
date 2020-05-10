package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("run")
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("./public/"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	//main_ll.go文件
	mux.HandleFunc("/", index)

	//route.go
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/authenticate", authenticate)

	//route_thread.go
	mux.HandleFunc("/thread/new", newThread)
	mux.HandleFunc("/thread/create", createThread)
	mux.HandleFunc("/thread/post", postThread)
	mux.HandleFunc("/thread/read", readThread)

	server := &http.Server{
		Addr:           "0.0.0.0:80",
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()

}
