package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

)

func main() {

	log.Print("QuotaSetter is starting...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/set", SetQuotas )
	router.HandleFunc("/get/{todoId}", GetQuotas)

	log.Fatal(http.ListenAndServe(":8080", router))

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to QuotaSetter!\n")
}

func SetQuotas(w http.ResponseWriter, r *http.Request) {
	var todo Todo//map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &todo)
	if err != nil { //输入参数解析出错
		log.Print("Can not decode data: %v\n", err)
		if error := json.NewEncoder(w).Encode(Result{Code: 10, IsSuccess: false, Description: err.Error()}); err != nil {
			panic(error)
		}
	}
	DoSetQuotas(w,todo)
}

func GetQuotas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintf(w, "Todo show: %s\n", todoId)
}

func DoSetQuotas(w http.ResponseWriter,todo Todo) error{
	log.Print("DoSetQuota,param:",todo)
	if(todo.Max_files!=0){
		cmd := exec.Command("setfattr ", "-n ceph.quota.max_files", "-v "+string(todo.Max_files),todo.Path)
		err := cmd.Run()
		if err != nil {//设置max_files出错 #11
			log.Print("Do cmd failed: ", err)
			if error := json.NewEncoder(w).Encode(Result{Code:11,IsSuccess:false,Description:err.Error()}); error != nil {
				log.Print(error)
				}
			return err
			}
		log.Print("Do cmd(set max_files) success.\n")
		}
	if(todo.Max_bytes!=0){
		cmd := exec.Command("setfattr ", "-n ceph.quota.max_bytes ", "-v "+string(todo.Max_bytes),todo.Path)
		err := cmd.Run()
		if err != nil {//设置max_bytes出错 #12
			log.Print("Do cmd failed: ", err)
			if error := json.NewEncoder(w).Encode(Result{Code:12,IsSuccess:false,Description:err.Error()}); error != nil {
				log.Print(error)
			}
			return err
		}
		log.Print("Do cmd(set max_bytes) success.\n")
	}
	if error := json.NewEncoder(w).Encode(Result{Code:20,IsSuccess:true,Description:"Set quotas success."}); error != nil {
		log.Print(error)
	}
	return nil
}