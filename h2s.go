package main

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/http2"

	"golang.org/x/net/http2/h2c"
)

/*
func main() {
	h2s := &http2.Server{}
	addr := "localhost:8080"
	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(http.HandlerFunc(requestHandler), h2s),
	}
	fmt.Printf("Listening %s...\n", addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("ERROR")
	}

}
func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		fmt.Println("got get method")
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}
*/
// var conns = make(map[string]net.Conn)

// func getConn(r *http.Request) net.Conn {
// 	return conns[r.RemoteAddr]
}
func main() {
	h2s := &http2.Server{}

	handler := http.HandlerFunc(requestHandle)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: h2c.NewHandler(handler, h2s),
	}

	fmt.Printf("Listening [localhost:8080]...\n")
	err := (server.ListenAndServe())
	if err != nil {
		fmt.Println(err)
	}
}
//======================================================================================

//======================================================================================
func requestHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/nchf-convergedcharging/v2/chargingdata" {
		fmt.Println("got it", r.Method)
		w.Header().Set("Response code", "201")
		fmt.Println(r.Body)
	}
	// conn := getConn(r)
	// fmt.Fprintln(conn, "Hello from tcp server")

	// if r.Method==post
	// fmt.Println(r.Header)

	fmt.Fprintf(w, "Hello, %v", r.URL.Path))
}
