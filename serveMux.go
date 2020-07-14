package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func main() {
	rot := Route{
		"HandleTimeout",
		strings.ToUpper("GET"),
		"/abc/def/ghi",
		HandleMessage,
	}

	var handlers http.HandlerFunc
	handlers = rot.HandlerFunc
	handlers = wraper(handlers)

	r := mux.NewRouter()
	r.Methods(rot.Method).Path(rot.Pattern).Name(rot.Name).Handler(wraper(handlers))

	//r.HandleFunc("/abc", wraper(doabc))

	srv := &http.Server{
		Handler: r,
		Addr:    "localhost:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
func HandleMessage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path
	fmt.Println("======Endpoint====== handler =============== called")
	io.WriteString(w, "=========this is abc route with doabc handler============")
}

func wraper(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Code for the middleware...
		fmt.Println("this is wrapper or middleware before")
		fmt.Fprint(w, "this is wrapper or middleware before")
		h.ServeHTTP(w, r)
		//doabc(w, r)
		fmt.Fprint(w, "this is wrapper or middleware after")
	})
}
