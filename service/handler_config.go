package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bytedance/douyincloud-configcenter-sdk-go/base"
)

// ConfigGetHandler 根据key获取一个配置，当key不存在时返回error
func ConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	//Count("/config/get")
	data := make(map[string]string)
	sdkClient, _ := base.Start()

	var keyObj map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &keyObj)
	value, err := sdkClient.Get(keyObj["key"].(string))
	data["key"] = keyObj["key"].(string)
	if err != nil {
		data["value"] = err.Error()
	} else {
		data["value"] = value
	}
	msg, _ := json.Marshal(data)
	w.Write(msg)
}

// ConfigRefreshHandler 刷新配置中心
func ConfigRefreshHandler(w http.ResponseWriter, r *http.Request) {
	//Count("/config/refresh")
	sdkClient, _ := base.Start()
	err := sdkClient.RefreshConfig()
	data := make(map[string]error)
	data["err"] = err
	msg, _ := json.Marshal(data)
	w.Write(msg)
}
