//package server

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/gorilla/mux"
)

var cid = 1000

var maap = make(map[int]int)

var ch = make(chan int, len(maap))

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func del(iid int) {
	fmt.Println("NUMBER OF ROUTINE RUNING IN DEL", runtime.NumGoroutine())
	tmr := time.NewTimer(time.Duration(30) * time.Second)
	for {
		select {
		case <-tmr.C:
			fmt.Println("DELETING THE IID==>>", iid)
			delete(maap, iid)
			fmt.Println("IID DELETED==>>", iid)
			tmr.Stop()
			return

		case val := <-ch:
			fmt.Println("in **********val := <-ch************ ", val, iid)
			if val == iid {
				tmr.Stop()
				tmr.Reset(time.Duration(20) * time.Second)
				fmt.Println("timer reset done")
			}
			continue

		}
	}
}

type Routes []Route

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
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

}
func broadcast(paramInt int) {
	for i := 1; i <= len(maap); i++ {
		ch <- paramInt
	}
	fmt.Println("BROADCAST DONE")
	return
}

///////////////////////////////////////////if url comes with variable id/////////////////////////////////////////////////////////////////
func handleuuid(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*************************NUMBER OF ROUTINE RUNING IN HNDLnoURL****************", runtime.NumGoroutine())
	fmt.Println("handleuuid")
	params := mux.Vars(r)["ID"]
	paramInt, _ := strconv.Atoi(params)

	go broadcast(paramInt)

	if len(maap) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
	ii := 0
	for index := range maap {

		if paramInt == index {
			newPath := "/nchf-convergedcharging/v2/chargingdata/" + params
			w.Header().Set("request-uri", newPath)
			hndl(w, r)
			break
		} else if ii == (len(maap) - 1) {
			w.WriteHeader(http.StatusNoContent)
		}
		ii++
	}
}

//handleMessage handles if url comes without url id
func handleMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NUMBER OF ROUTINE RUNING IN HNDLURL", runtime.NumGoroutine())
	fmt.Println("entered handlemessage")
	uuid := cid
	maap[uuid] = uuid
	go del(uuid)

	fmt.Println("===>", maap)
	temp := strconv.Itoa(uuid) //int to string
	fmt.Println("temp=", temp)
	newPath := "/nchf-convergedcharging/v2/chargingdata/" + temp
	fmt.Println(newPath)
	w.Header().Set("request-uri", newPath)
	hndl(w, r)
	cid++
	//count = 2

	return

}

//================================================================

func main() {
	fmt.Println("NUMBER OF ROUTINE RUNING", runtime.NumGoroutine())
	router := NewRouter()
	listnerPort := fmt.Sprintf(":%d", 8080)

	h2server := &http2.Server{}
	server := &http.Server{
		Addr:        listnerPort,
		Handler:     h2c.NewHandler(router, h2server), //here router is nothing but serveMux which itself is Handler, as it is having method serveHttp attached to it
		IdleTimeout: 100 * time.Second,
		//ConnState:   listenit,
	}
	log.Println("Started listening on port 8080")
	fmt.Println("============")
	if err := server.ListenAndServe(); err != nil {
		log.Println("failed to listen server : ", err)
		os.Exit(1)

	}
	fmt.Println("======||======")

	fmt.Println("end statement")

}
