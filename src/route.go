package main

import (
	"io"
	"net/http"
	"regexp"
	"time"
	"util"
)

type WebController struct{
	Function func(http.ResponseWriter,* http.Request)
	Method string
	Pattern string
}

var mux []WebController

func  init(){
	mux = append(mux,WebController{post,"POST","^/"})
	mux = append(mux,WebController{get,"GET","^/"})
}

type httpHandler struct{}

func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t := time.Now()

	for _, webController := range mux {

		if m, _ := regexp.MatchString(webController.Pattern, r.URL.Path); m {

			if r.Method == webController.Method {
				webController.Function(w, r)
				go util.WriteLog(r, t, "match", webController.Pattern)
				return
			}
		}
	}

	go util.WriteLog(r, t, "un match", "")

	io.WriteString(w, "")
	return
}