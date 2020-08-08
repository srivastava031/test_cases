//package server

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/gorilla/mux"
)

var cid = 1000
var count = 1

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//var m = make(map[string]time.Time)
var uid []int

type Routes []Route

var doneCh = make(chan struct{})

var routes = Routes{
	Route{
		"HandleTimeout",
		strings.ToUpper("Post"),
		"/nchf-convergedcharging/v2/chargingdata/",
		handleMessage,
	},
	{
		"HandleTimeout",
		strings.ToUpper("Post"),
		"/nchf-convergedcharging/v2/chargingdata/{ID}",
		handleuuid,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		//handler = Logger(handler, route.Name) //decorator

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler) //this is calling function
	}

	return router
}

//===============================================================================================
//handler
//================================================================================================
func hndl(w http.ResponseWriter, r *http.Request) {
	content := make(map[string]interface{})
	reqBody, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	re := bytes.NewReader(reqBody)
	if err := json.NewDecoder(re).Decode(&content); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//log.Println("error bad request :", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

}

//if url comes with variable id
func handleuuid(w http.ResponseWriter, r *http.Request) {
	if count > 1 {
		go cancel()
	}
	fmt.Println("handleuuid")
	params := mux.Vars(r)["ID"]
	paramInt, _ := strconv.Atoi(params)
	if uid == nil {
		w.WriteHeader(http.StatusNoContent)
	}
	for index, value := range uid {
		if paramInt == value {
			newPath := "/nchf-convergedcharging/v2/chargingdata/" + params
			w.Header().Set("request-uri", newPath)
			hndl(w, r)
			fmt.Println("got in if loop")
			break
		} else if index == (len(uid) - 1) {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

//if url comes without url id
func handleMessage(w http.ResponseWriter, r *http.Request) {
	if count > 1 {
		go cancel()
	}
	//t:=time.Now()
	uuid := cid
	uid = append(uid, uuid)
	fmt.Println(uid)
	temp := strconv.Itoa(uuid) //int to string
	fmt.Println("temp=", temp)
	newPath := "/nchf-convergedcharging/v2/chargingdata/" + temp
	fmt.Println(newPath)
	w.Header().Set("request-uri", newPath)
	hndl(w, r)
	cid++
	count = 2

	return

}

//================================================================

func main() {
	// ctx:=context.Background()
	// ctx, cancel:=context.WithCancel(ctx)
	router := NewRouter()
	listnerPort := fmt.Sprintf(":%d", 8080)

	h2server := &http2.Server{}
	server := &http.Server{
		Addr:        listnerPort,
		Handler:     h2c.NewHandler(router, h2server), //here router is nothing but serveMux which itself is Handler, as it is having method serveHttp attached to it
		IdleTimeout: 100 * time.Second,
		ConnState:   connStateListener,
	}
	log.Println("Started listening on port 8080")
	fmt.Println("============")
	if err := server.ListenAndServe(); err != nil {
		log.Println("failed to listen server : ", err)
		os.Exit(1)

	}
	fmt.Println("============")
	// notify := make(chan error)

}

func connStateListener(c net.Conn, cs http.ConnState) {
	switch cs {
	case http.StateIdle:
		go idletime()
		fmt.Println("idletime")
	}

}

func idletime() {
	fmt.Println("entered idletime")
	select {
	case <-time.After(20 * time.Second):
		uid = nil
		count = 1
		cid = 1000
		fmt.Println(uid)
	case <-doneCh:

		fmt.Println("encountered done")
		return
	}

}

func cancel() {
	doneCh <- struct{}{}
}
