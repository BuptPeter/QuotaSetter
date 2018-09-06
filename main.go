package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io/ioutil"

	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true                    // 显示完整时间
	customFormatter.TimestampFormat = "2006-01-02 15:04:05" // 时间格式
	customFormatter.DisableTimestamp = false                // 禁止显示时间
	customFormatter.DisableColors = false                   // 禁止颜色显示
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
func main() {
	log.Info("QuotaSetter is starting...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/set", SetQuotas)
	router.HandleFunc("/get", GetQuotas)
	log.Fatal(http.ListenAndServe(":8081", router))

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to QuotaSetter!\n")
}

func SetQuotas(w http.ResponseWriter, r *http.Request) {
	log.Info("handle /set request.")
	todo := Todo{Max_bytes: -1, Max_files: -1} //map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &todo)
	if err != nil { //输入参数解析出错
		if error := json.NewEncoder(w).Encode(SetResult{Code: 10, IsSuccess: false, Description: err.Error()}); error != nil {
			log.Error(error)
		}
		log.Error("Can not decode data: \n", err)
		return
	}
	DoSetQuotas(w, todo)
}
func GetQuotas(w http.ResponseWriter, r *http.Request) {
	log.Info("handle /get request.")
	var todo Todo
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &todo)
	if err != nil { //输入参数解析出错
		if error := json.NewEncoder(w).Encode(GetResult{Code: 10, IsSuccess: false, Description: err.Error()}); error != nil {
			log.Error(error)
		}
		log.Error("Can not decode data: \n", err)
		return
	}
	DoGetQuotas(w, todo)
}
func DoGetQuotas(w http.ResponseWriter, todo Todo) {
	log.Info("DoGetQuota,param:", todo.Path, (todo.Max_bytes), (todo.Max_files))
	if isDirExists(todo.Path) == false { //目标目录不存在 #11
		if error := json.NewEncoder(w).Encode(GetResult{Code: 11, IsSuccess: false, Description: "target path is not exist."}); error != nil {
			log.Error(error)
		}
		log.Error("target path is not exist.")
		return
	}
	var stdErr0, stdOut0 bytes.Buffer
	cmd0 := exec.Command("/bin/bash", "-c", "getfattr -n ceph.quota.max_files "+todo.Path)
	cmd0.Stderr = &stdErr0
	cmd0.Stdout = &stdOut0
	err := cmd0.Run()
	if err != nil { //获取max_files出错 #15
		if error := json.NewEncoder(w).Encode(GetResult{Code: 15, IsSuccess: false, Description: err.Error() + stdErr0.String()}); error != nil {
			log.Error(error)
		}
		log.Error("Do cmd failed: \n", err.Error()+stdErr0.String())
		return
	}
	log.Info("Do cmd(get max_files) success.\n")
	var stdErr1, stdOut1 bytes.Buffer
	cmd1 := exec.Command("/bin/bash", "-c", "getfattr  -n ceph.quota.max_bytes "+todo.Path)
	cmd1.Stderr = &stdErr1
	cmd1.Stdout = &stdOut1
	err = cmd1.Run()
	if err != nil { //获取max_bytes 出错 #16
		if error := json.NewEncoder(w).Encode(GetResult{Code: 16, IsSuccess: false, Description: err.Error() + stdErr1.String()}); error != nil {
			log.Error(error)
		}
		log.Error("Do cmd failed: \n", err.Error()+stdErr1.String())
		return
	}
	log.Info("Do cmd(get max_files) success.\n")
	max_bytes, _ := strconv.Atoi(strings.Split(stdOut1.String(), "\"")[1])
	max_files, _ := strconv.Atoi(strings.Split(stdOut0.String(), "\"")[1])
	if error := json.NewEncoder(w).Encode(GetResult{Code: 21, IsSuccess: true, Description: "Get quotas success.", Max_bytes: max_bytes, Max_files: max_files}); error != nil {
		log.Error(error)
	}
}
func DoSetQuotas(w http.ResponseWriter, todo Todo) {
	log.Info("DoSetQuota,param:", todo)
	if isDirExists(todo.Path) == false { //目标目录不存在 #11
		if error := json.NewEncoder(w).Encode(SetResult{Code: 11, IsSuccess: false, Description: "target path is not exist."}); error != nil {
			log.Error(error)
		}
		log.Error("Target path is not exist.")
		return
	}
	if (todo.Max_files) < 0 || (todo.Max_bytes) < 0 {//缺少参数或参数无效 #12
		if error := json.NewEncoder(w).Encode(SetResult{Code: 12, IsSuccess: false, Description: "Invalid argument."}); error != nil {
			log.Error(error)
		}
		log.Error("Missing parameters or Invalid parameters.")
		return
	}
	if todo.Max_files >=0 {
		var out bytes.Buffer
		//cmd0 := exec.Command("/bin/bash", "-c","\"setfattr ", "-n ceph.quota.max_files", "-v",string(todo.Max_files), todo.Path+"\"")
		cmd := "setfattr  -n ceph.quota.max_files -v " + strconv.Itoa(todo.Max_files) + " " + todo.Path
		cmd0 := exec.Command("/bin/bash", "-c", cmd)
		cmd0.Stderr = &out
		err := cmd0.Run()
		if err != nil { //设置max_files出错 #13
			if error := json.NewEncoder(w).Encode(SetResult{Code: 13, IsSuccess: false, Description: out.String()}); error != nil {
				log.Error(error)
			}
			log.Error("Do cmd ("+cmd+") failed:\n ", out.String())
			return
		}
		log.Info("Do cmd(set max_files) success. \n")
	}
	if todo.Max_bytes >=0 {
		var out bytes.Buffer
		//cmd1 := exec.Command("/bin/bash", "-c","\"setfattr ", "-n ceph.quota.max_bytes ", "-v "+string(todo.Max_bytes), todo.Path+"\"")
		cmd := "setfattr -n ceph.quota.max_bytes -v " + strconv.Itoa(todo.Max_bytes) + " " + todo.Path
		cmd1 := exec.Command("/bin/bash", "-c", cmd)
		cmd1.Stderr = &out
		err := cmd1.Run()
		if err != nil { //设置max_bytes出错 #14
			if error := json.NewEncoder(w).Encode(SetResult{Code: 14, IsSuccess: false, Description: out.String()}); error != nil {
				log.Error(error)
			}
			log.Error("Do cmd ("+cmd+") failed: \n", out.String())
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
