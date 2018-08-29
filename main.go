package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true                      // 显示完整时间
	customFormatter.TimestampFormat = "[2006-01-02 15:04:05]" // 时间格式
	customFormatter.DisableTimestamp = false                  // 禁止显示时间
	customFormatter.DisableColors = false                     // 禁止颜色显示
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.Info("QuotaSetter is starting...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/set", SetQuotas)
	router.HandleFunc("/get/{TargetPath}", GetQuotas)
	log.Fatal(http.ListenAndServe(":8080", router))

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to QuotaSetter!\n")
}

func SetQuotas(w http.ResponseWriter, r *http.Request) {
	var todo Todo //map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &todo)
	if err != nil { //输入参数解析出错
		if error := json.NewEncoder(w).Encode(SetResult{Code: 10, IsSuccess: false, Description: err.Error()}); err != nil {
			log.Error(error)
		}
		log.Error("Can not decode data: %v\n", err)
		return
	}
	DoSetQuotas(w, todo)
}

func GetQuotas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	TargetPath := vars["TargetPath"]
	log.Info("GetQuota,param:", TargetPath)
	fmt.Fprintf(w, "Todo show: %s\n", TargetPath)
	if isDirExists(TargetPath) == false { //目标目录不存在 #11
		if error := json.NewEncoder(w).Encode(SetResult{Code: 11, IsSuccess: false, Description: "target path is not exist."}); error != nil {
			log.Error(error)
		}
		log.Error("target path is not exist.")
		return
	}
	cmd := exec.Command("getfattr  ", "-n ceph.quota.max_files", TargetPath)
	err := cmd.Run()
	if err != nil { //获取max_files出错 #14
		if error := json.NewEncoder(w).Encode(SetResult{Code: 14, IsSuccess: false, Description: err.Error()}); error != nil {
			log.Error(error)
		}
		log.Error("Do cmd failed: ", err)
		return
	}
	log.Info("Do cmd(get max_files) success.\n")

	cmd = exec.Command("getfattr  ", "-n ceph.quota.max_bytes ", TargetPath)
	err = cmd.Run()
	if err != nil { //获取max_bytes 出错 #15
		if error := json.NewEncoder(w).Encode(SetResult{Code: 15, IsSuccess: false, Description: err.Error()}); error != nil {
			log.Error(error)
		}
		log.Error("Do cmd failed: ", err)
		return
	}
	log.Info("[INFO]:Do cmd(get max_files) success.\n")

}

func DoSetQuotas(w http.ResponseWriter, todo Todo) {
	log.Info("DoSetQuota,param:", todo)
	if isDirExists(todo.Path) == false { //目标目录不存在 #11
		if error := json.NewEncoder(w).Encode(SetResult{Code: 11, IsSuccess: false, Description: "target path is not exist."}); error != nil {
			log.Error(error)
		}
		log.Error("target path is not exist.")
		return
	}
	if todo.Max_files < 0 || todo.Max_bytes < 0 {
		if error := json.NewEncoder(w).Encode(SetResult{Code: 11, IsSuccess: false, Description: "Invalid argument."}); error != nil {
			log.Error(error)
		}
		log.Error("Invalid argument.")
		return
	}
	//cmd := exec.Command("cmd", "/C")
	if todo.Max_files >= 0 {
		cmd := exec.Command("setfattr ", "-n ceph.quota.max_files", "-v "+string(todo.Max_files), todo.Path)
		err := cmd.Run()
		if err != nil { //设置max_files出错 #12
			if error := json.NewEncoder(w).Encode(SetResult{Code: 12, IsSuccess: false, Description: err.Error()}); error != nil {
				log.Error(error)
			}
			log.Error("Do cmd failed: ", err)
			return
		}
		log.Info("Do cmd(set max_files) success.\n")
	}
	//cmd1 := exec.Command("cmd", "/C")
	if todo.Max_bytes >= 0 {
		cmd := exec.Command("setfattr ", "-n ceph.quota.max_bytes ", "-v "+string(todo.Max_bytes), todo.Path)
		err := cmd.Run()
		if err != nil { //设置max_bytes出错 #13

			if error := json.NewEncoder(w).Encode(SetResult{Code: 13, IsSuccess: false, Description: err.Error()}); error != nil {
				log.Error(error)
			}
			log.Error("Do cmd failed: ", err)
			return
		}
		log.Info("Do cmd(set max_bytes) success.\n")
	}
	if error := json.NewEncoder(w).Encode(SetResult{Code: 20, IsSuccess: true, Description: "Set quotas success."}); error != nil {
		log.Error(error)
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
