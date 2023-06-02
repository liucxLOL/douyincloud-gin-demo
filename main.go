/*
Copyright (year) Bytedance Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"douyincloud-gin-demo/db/mongodb"
	"douyincloud-gin-demo/db/mysql"
	"douyincloud-gin-demo/db/redis"
	"douyincloud-gin-demo/service"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	mysql.InitMysql()
	redis.InitRedis()
	mongodb.InitMongoDB()

	http.HandleFunc("/v1/ping", service.PingHandler) //火山校验

	http.HandleFunc("/api/getQuestionnaireList", service.SelectQuestionnaireList)
	http.HandleFunc("/api/getQuestionnaireInfo", service.GetQuestionnaireInfo)
	http.HandleFunc("/api/createQuestionnaire", service.CreateQuestionnaireInfo)
	http.HandleFunc("/mysql/create_lock_table", service.MysqlCreateLockTable)
	http.HandleFunc("/mysql/update", service.MysqlUpdate)
	http.HandleFunc("/mysql/update_counts", service.MysqlUpdateCounts)
	http.HandleFunc("/mysql/delete", service.MysqlDelete)
	http.HandleFunc("/mysql/delete_rollback", service.MysqlDeleteRollback)

	listenPort := ":8000"
	if listenPort == "" {
		log.Fatal("failed to load _FAAS_RUNTIME_PORT")
	}

	// 部署状态控制
	if os.Getenv("deploy") == "fail" {
		panic("crash")
		listenPort = ""
		fmt.Println("fail")
	}

	fmt.Println("http ListenAndServe ", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
