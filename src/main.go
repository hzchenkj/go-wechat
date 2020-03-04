package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"wx"
)

const (
	port  = 8080
	token = "7UaMXnuvFD8SoBMW"
)

func get(w http.ResponseWriter, r *http.Request) {
	log.Println("get method 第一次认证用")
	client, err := wx.NewClient(r, w, token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if len(client.Query.EchoStr) > 0 {
		w.Write([]byte(client.Query.EchoStr))
		return
	}
	w.WriteHeader(http.StatusForbidden)
	return
}

func post(w http.ResponseWriter, r *http.Request) {
	log.Println("post method")
	client, err := wx.NewClient(r, w, token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	client.Run()
	return
}

func main() {
	// http server
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        &httpHandler{},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	log.Println(fmt.Sprintf("Listen: %d", port))
	log.Fatal(server.ListenAndServe())
}
