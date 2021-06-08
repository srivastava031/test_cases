package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"encoding/json"
)
type employee struct {
	Name string `json:"name"`
	EmployeeId string `json:"employeeId"`
}
func createClient() (http.Client) {
	var client http.Client
	// http2 will use its parameter for client configuration but will not use it as transport
	client = http.Client{
		Timeout: 5*time.Second,
	}
	return client
}
//SendRequest sends the request to the URL given
func sendRequest(client http.Client, method string, url string, header map[string]string, body []byte) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
	}
	if header != nil {
		for key, value := range header {
			req.Header[key] = []string{value}
		}
	}
	response, err := client.Do(req)
	if response == nil {
		fmt.Println("no response recieved!!")
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(response.Status)
	fmt.Println(string(body))

}

func main() {
	client:=createClient()
	method:=flag.String("M", "GET","method")
	pathParameter := flag.String("P", "", "a string")
	name:=flag.String("N", "","employee name")
	id := flag.String("I", "", "employee id")
	headerMap:=map[string]string{"Content-Type": "application/json; charset=utf-8"}
	flag.Parse()

	body := map[string]string{*name:*id}

	var serialisedBody []byte
	
	if *method=="POST" {
		serialisedBody,_=json.Marshal(body)
	}
	sendRequest(client, (*method), "http://127.0.0.1:8090/employee/"+(*pathParameter), headerMap,serialisedBody)
	}

