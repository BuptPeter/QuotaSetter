package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

)

func main() {

	log.Print("QuotaSetter is starting...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/set", SetQuotas )
	router.HandleFunc("/get/{todoPath}", GetQuotas)

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
		if error := json.NewEncoder(w).Encode(SetResult{Code: 10, IsSuccess: false, Description: err.Error()}); err != nil {
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

func DoSetQuotas(w http.ResponseWriter,todo Todo) {
	log.Print("DoSetQuota,param:",todo)
	if(isDirExists(todo.Path)==false){//目标目录不存在 #11
		log.Print("[ERROR]:target path is not exist.")
		if error := json.NewEncoder(w).Encode(SetResult{Code:11,IsSuccess:false,Description:"target path is not exist."}); error != nil {
			log.Print("[ERROR]:",error)
		}
		return
	}
	if(todo.Max_files!=0){
		cmd := exec.Command("setfattr ", "-n ceph.quota.max_files", "-v "+string(todo.Max_files),todo.Path)
		err := cmd.Run()
		if err != nil {//设置max_files出错 #12
			log.Print("[ERROR]:Do cmd failed: ", err)
			if error := json.NewEncoder(w).Encode(SetResult{Code:12,IsSuccess:false,Description:err.Error()}); error != nil {
				log.Print("[ERROR]:",error)
				}
			return
			}
		log.Print("[INFO]:Do cmd(set max_files) success.\n")
		}
	if(todo.Max_bytes!=0){
		cmd := exec.Command("setfattr ", "-n ceph.quota.max_bytes ", "-v "+string(todo.Max_bytes),todo.Path)
		err := cmd.Run()
		if err != nil {//设置max_bytes出错 #13
			log.Print("[ERROR]:Do cmd failed: ", err)
			if error := json.NewEncoder(w).Encode(SetResult{Code:13,IsSuccess:false,Description:err.Error()}); error != nil {
				log.Print("[ERROR]:",error)
			}
			return
		}
		log.Print("[INFO]:Do cmd(set max_bytes) success.\n")
	}
	if error := json.NewEncoder(w).Encode(SetResult{Code:20,IsSuccess:true,Description:"Set quotas success."}); error != nil {
		log.Print("[ERROR]:",error)
	}
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("not reached")
}