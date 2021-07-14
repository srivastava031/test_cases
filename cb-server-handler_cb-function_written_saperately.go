package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strconv"
	"reflect"

)

var collection *gocb.Collection
var bucket *gocb.Bucket
var cluster *gocb.Cluster
var router = mux.NewRouter().StrictSlash(true)

type ProjectDetails struct {
	Language string `json:"language,omitempty"`
	Platform string `json:"platform,omitempty"`
}

type Employee struct {
	Name        string                    `json:"name,omitempty"`
	EmpId       string                    `json:"emp_id,omitempty"`
	Salary      string                    `json:"salary,omitempty"`
	Address     string                    `json:"address,omitempty"`
	PhoneNumber string                    `json:"phone_number,omitempty"`
	Projects    map[string]ProjectDetails `json:"projects,omitempty"`
}

func cbGEToperation(id,field ,url string)(int, interface{}, error){
	if url=="/employee/4?field="{
	v, err := collection.Get(id, &gocb.GetOptions{})//
	if err!=nil{
		return 417, nil, err
	}
	var value interface{}
	err = v.Content(&value)
	if err!=nil{
		return 417, nil, err
	}else{
		return 200, value, nil
		}
	 }else{
		var ops []gocb.LookupInSpec
			ops = []gocb.LookupInSpec{
				gocb.GetSpec(field, &gocb.GetSpecOptions{}),
				}
			getResult, err := collection.LookupIn(id, ops, &gocb.LookupInOptions{})
			if err != nil {
				return 417, nil, err
			}
			var value interface{}
			err = getResult.ContentAt(0, &value)
			if err != nil {
				return 417, nil, err
	 }
	 return 200, value, nil
	}
	return 0, nil, nil
}

func cbPOSToperation(id string, body []byte)(int, error){
	employee := Employee{}
	err := json.Unmarshal(body,&employee)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return 417, err
	}

	if employee.EmpId != id{
		return 401, err
	}
	_, err = collection.Insert(id, employee ,&gocb.InsertOptions{})
	if err!=nil{
		return 417, err
	}	
	return 201, nil
}

func cbPUToperation(id string, body []byte)(int,error){
		updatedEmployee := make(map[string]interface{})
		err := json.Unmarshal(body,&updatedEmployee)
		if err != nil {
			return 203,err
		}
		var updateDoc []gocb.MutateInSpec
		var tempVar string
		for k,v := range updatedEmployee{
			if reflect.TypeOf(v).Kind()!=reflect.Map{
		 	updateDoc = append(updateDoc,gocb.UpsertSpec(k, v, &gocb.UpsertSpecOptions{}))
			}
			 if reflect.ValueOf(v).Kind()==reflect.Map{
				tempVar=fmt.Sprintf("%v",k)
				vmap:=v.(map[string]interface{})
				for key,val := range vmap{
					tempVar:=fmt.Sprintf("%v.%v",tempVar,key)
					if reflect.TypeOf(val).Kind()!=reflect.Map{
					updateDoc = append(updateDoc,gocb.UpsertSpec(tempVar, val, &gocb.UpsertSpecOptions{}))
					}
					if reflect.TypeOf(val).Kind()==reflect.Map{
						vmap:=val.(map[string]interface {})
						for key,val := range vmap{
							tempVar:=fmt.Sprintf("%v.%v",tempVar,key)
							updateDoc = append(updateDoc,gocb.UpsertSpec(tempVar, val, &gocb.UpsertSpecOptions{}))
						}
					}
				}
			 }
		}
		_, err = collection.MutateIn(id, updateDoc, &gocb.MutateInOptions{})
		if err != nil {
			return 417, err
		}
		return 202, nil

	}

func cbDELETEoperation(id string)(int,error){
	if id!="all"{
		_, err := collection.Remove(id,&gocb.RemoveOptions{})
		if err!=nil{
			return 417, err
			}	
	}else if id=="all"{
		var itemss []gocb.BulkOp
		for i := 1; i < 10; i++ {
			itemss = append(itemss, &gocb.RemoveOp{ID: strconv.Itoa(i)})
		}
		err := collection.Do(itemss,&gocb.BulkOpOptions{})
		
			if err != nil {
				return 417, err
			}
	}
	return 202,nil
}
// will return substring field details 
func getEmployeeDetail(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	field := r.FormValue("field")// form value is used because it will place empty string where query parameter is nil
	id := mux.Vars(r)["id"]

	u, err := router.Get("GETemployee").URL("id", id, "field", field)//used to build back the url from the router name.
	if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
	url:=u.String()

	statusCode,value,err:=cbGEToperation(id,field,url)	
	if err!=nil{
		http.Error(w, err.Error(), statusCode)
	}else{
		json.NewEncoder(w).Encode(value)
		w.WriteHeader(statusCode)
	}
}

// Add new employee
func createEmployeeDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	param := params["id"]
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	statusCode,err:=cbPOSToperation(param,body)
	if err!=nil{
		http.Error(w, err.Error(), statusCode)
	}else{
	w.WriteHeader(statusCode)
	}
}

//Update the name of employee
func update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param:= params["id"]
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	statusCode,err:=cbPUToperation(param, body)
	if err!=nil{
		http.Error(w, err.Error(), statusCode)
	}else{
	w.WriteHeader(statusCode)
	}
}

// Delete employee
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	param:= params["id"]

	statusCode,err:=cbDELETEoperation(param)
	if err!=nil{
		http.Error(w, err.Error(), statusCode)
		}else{
			w.WriteHeader(statusCode)
		}
	}


// Main function
func main() {
	var err error
	cluster, err = gocb.Connect(
		"localhost",
		gocb.ClusterOptions{
			Username: "Administrator",
			Password: "password",
		})
	if err != nil {
		panic(err)
	}
	bucket = cluster.Bucket("employee")
	collection = bucket.DefaultCollection()
	router.HandleFunc("/employee/{id}", getEmployeeDetail).Queries("field", "{field}").Name("GETemployee").Methods("GET")
	router.HandleFunc("/employee/{id}", getEmployeeDetail).Methods("GET")
	router.HandleFunc("/employee/{id}", createEmployeeDetail).Methods("POST")
	router.HandleFunc("/employee/{id}", deleteEmployee).Methods("DELETE")
	router.HandleFunc("/employee/{id}", update).Methods("PUT")
	fmt.Println("started server...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8090", router))
}
































//https://docs.couchbase.com/go-sdk/current/howtos/subdocument-operations.html

	// var items []gocb.BulkOp
	// var itemsGet []gocb.BulkOp


	// if field=="name"{
	// 	mops := []gocb.MutateInSpec{
	// 		gocb.UpsertSpec("name", newName.Name, &gocb.UpsertSpecOptions{}),
	// 	}
	// }else if field=="salary"{
	// 	ops = []gocb.LookupInSpec{
	// 		gocb.GetSpec("salary", &gocb.GetSpecOptions{}),
	// 		}

	// func getAll(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("------getAll---------------------------")
	// 	cluster, err := gocb.Connect(
	// 		"localhost",
	// 		gocb.ClusterOptions{
	// 			Username: "Administrator",
	// 			Password: "password",
	// 		})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	// bucket := cluster.Bucket("employee")
	// 	// collection := bucket.DefaultCollection()
	// 	w.Header().Set("Content-Type", "application/json")
	
	// 	results, err := cluster.Query("SELECT * FROM `employee`", nil)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	//g:=gocb.GetAllQueryIndexesOptions{}
	// 	//var greeting interface{}
	// 	m:=make(map[string]map[string]interface{})
	// 	mm:=make(map[interface{}]interface{})
	
	// 	for results.Next() {
	// 		err := results.Row(&m)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		eid:=(m["employee"]["emp_id"])
	// 		mm[eid]=eid
	// 	}
	// 	fmt.Println(mm)
	
	// 	// always check for errors after iterating
	// 	err = results.Err()
	// 	if err != nil {
	// 		panic(err)
	
	// 	// qi,err:=ind.GetAllIndexes("employee",&gocb.GetAllQueryIndexesOptions{})
	// 	// if err!=nil{
	// 	// 	json.NewEncoder(w).Encode(err)
	// 	// 	return
	// 	// 	}
	// 	// 	fmt.Println(qi)
	// 	// 	json.NewEncoder(w).Encode(qi)
	// 	}
	// }
