package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"reflect"
	"strings"
	"strconv"

)

var collection *gocb.Collection
var bucket *gocb.Bucket
var cluster *gocb.Cluster

// type ProjectDetails struct {
// 	Language string `json:"language,omitempty"`
// 	Platform string `json:"platform,omitempty"`
// }

type Employee struct {
	Name        string                 `json:"name,omitempty"`
	EmpId       string                 `json:"emp_id,omitempty"`
	Salary      string                 `json:"salary,omitempty"`
	Address     string                 `json:"address,omitempty"`
	PhoneNumber string                 `json:"phone_number,omitempty"`
	Projects    map[string]interface{} `json:"projects,omitempty"`
}

func cbGEToperation(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r) // Gets params
	param := params["id"]
	field:=params["field"]
	if field !="all"{
	emp:= Employee{}
	v := reflect.ValueOf(emp)
	typeOfS := v.Type()
	for i := 0; i< v.NumField(); i++ {
		if strings.ToLower(field)==strings.ToLower(typeOfS.Field(i).Name){
			break
		}else if i==v.NumField()-1{
			json.NewEncoder(w).Encode("no matching field found")
			return
		}
	}
	var ops []gocb.LookupInSpec
	ops = []gocb.LookupInSpec{
		gocb.GetSpec(field, &gocb.GetSpecOptions{}),
		}
	getResult, err := collection.LookupIn(param, ops, &gocb.LookupInOptions{})
	if err != nil {
		json.NewEncoder(w).Encode("document not found")
		return
	}
	var value interface{}
	err = getResult.ContentAt(0, &value)
	if err != nil {
		json.NewEncoder(w).Encode("document not found")
		return
	}
	
	json.NewEncoder(w).Encode(value)
	w.WriteHeader(http.StatusOK)
 }else if field =="all"{
	v, err := collection.Get(params["id"], &gocb.GetOptions{})//
	if err!=nil{
		json.NewEncoder(w).Encode("document not found")
		return
	}
	var value interface{}
	err = v.Content(&value)
	if err!=nil{
		json.NewEncoder(w).Encode(err)
	}else{
		json.NewEncoder(w).Encode(value)
		w.WriteHeader(http.StatusOK)
		}
	}else{
		json.NewEncoder(w).Encode("not valid path parameter")
	}
}

func cbPOSToperation(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r) // Gets params
	param := params["id"]
	employee := Employee{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body,&employee)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	if employee.EmpId != param{
		json.NewEncoder(w).Encode("employeeID not matching")
		return
	}
	_, err = collection.Insert(param, employee ,&gocb.InsertOptions{})
	if err!=nil{
		fmt.Println("ERROR: document exists")
		json.NewEncoder(w).Encode("ERROR: document exists")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
	//json.NewEncoder(w).Encode(v.Result.Cas())	
}

func cbPUToperation(w http.ResponseWriter, r *http.Request){
	//time.Sleep(30*time.Second)
	params := mux.Vars(r)
	param:=params["id"]

	qim:=cluster.QueryIndexes()
	_ = qim.CreatePrimaryIndex("employee",&gocb.CreatePrimaryQueryIndexOptions{IgnoreIfExists:true})

	results, err := cluster.Query("SELECT * FROM `employee`", nil)
	if err != nil {
		panic(err)
	}

	recievedResource:=make(map[string]map[string]interface{})
	resourceMap:=make(map[interface{}]interface{})

	//var greeting interface{}
	for results.Next() {
		err := results.Row(&recievedResource)
		if err != nil {
			panic(err)
		}
		eid:=(recievedResource["employee"]["emp_id"])
		resourceMap[eid]=eid
	}

	_,ok:=resourceMap[param]
	if !ok{
		json.NewEncoder(w).Encode("resource for given id not found ")
		return
	}
	// always check for errors after iterating
	err = results.Err()
	if err != nil {
	
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	if params["field"]=="subdoc"{
		updatedEmployee := make(map[string]interface{})
		err = json.Unmarshal(body,&updatedEmployee)
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		v, ok := updatedEmployee["emp_id"]
		if ok && v!= param {
			json.NewEncoder(w).Encode("empId not matching")
			return
		}
		var updateDoc []gocb.MutateInSpec

		for k,v := range updatedEmployee{
		 	updateDoc = append(updateDoc,gocb.UpsertSpec(k, v, &gocb.UpsertSpecOptions{}))
		}
		_, err = collection.MutateIn(param, updateDoc, &gocb.MutateInOptions{})
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}else if params["field"]=="doc"{
		updatedEmployee := Employee{}
		err = json.Unmarshal(body,&updatedEmployee)
		if updatedEmployee.EmpId!="" && updatedEmployee.EmpId != param{
			json.NewEncoder(w).Encode("empId not matching")
			return
		}
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}
		_, err = collection.Upsert(param, updatedEmployee ,&gocb.UpsertOptions{})
		if err!=nil{
			fmt.Println(err)
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		//json.NewEncoder(w).Encode(employee)
		}	
	}

func cbDELETEoperation(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	if params["id"]!="all"{
		_, err := collection.Remove(params["id"],&gocb.RemoveOptions{})
		if err!=nil{
			json.NewEncoder(w).Encode("document not found")
			return
			}
			w.WriteHeader(http.StatusNoContent)	
	}else if params["id"]=="all"{
		var itemss []gocb.BulkOp
		for i := 1; i < 10; i++ {
			itemss = append(itemss, &gocb.RemoveOp{ID: strconv.Itoa(i)})
		}
		err := collection.Do(itemss,&gocb.BulkOpOptions{})
		
			if err != nil {
				json.NewEncoder(w).Encode("ERRROR PERFORMING BULK DELETE")
			}
	}
}


// will return substring field details 
func getEmployeeDetail(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	cbGEToperation(w,r)	
}

// Add new employee
func createEmployeeDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cbPOSToperation(w,r)
}

//Update the name of employee
func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cbPUToperation(w,r)
}

// Delete employee
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cbDELETEoperation(w,r)
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
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/employee/{field}/{id}", getEmployeeDetail).Methods("GET")
	r.HandleFunc("/employee/{id}", createEmployeeDetail).Methods("POST")
	r.HandleFunc("/employee/{id}", deleteEmployee).Methods("DELETE")
	r.HandleFunc("/employee/{field}/{id}", update).Methods("PUT")
	fmt.Println("started server...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8090", r))
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