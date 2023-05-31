package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type jsonLog struct {
	LogID  string  `json:"logID"`
	Age    int64   `json:"age"`
	Salary float64 `json:"salary"`
	Name   string  `json:"name"`
	Height float64 `json:"height"`
	Email  string  `json:"email"`
}

// 获取当前可用区
func getAzName(w http.ResponseWriter) (res string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Fprintf(w, "InterfaceAddrs error:%+v", err)
		return "err"
	}

	for _, address := range addrs {
		ipnet, _ := address.(*net.IPNet)
		if ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), "192.168.0.") {
			//fmt.Fprintf(w, "Ip:%s AzName:%s\n", ipnet.IP.String(), "可用区A")
			return "可用区A：" + ipnet.IP.String()
		} else if ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), "192.168.1.") {
			//fmt.Fprintf(w, "Ip:%s AzName:%s\n", ipnet.IP.String(), "可用区B")
			return "可用区B：" + ipnet.IP.String()
		}
	}
	return "可用区未知"
}

// TestHandler api/test接口，本函数为测试demo的重要入口，请按照注释修改
func TestHandler(w http.ResponseWriter, r *http.Request) {
	//Count("/api/test")
	//定义变量
	var azName = getAzName(w)

	//打印日志
	fmt.Println("CLOUD_ENV: ", os.Getenv("CLOUD_ENV"))
	fmt.Println("grey_debug_flag: ", os.Getenv("flag"))
	fmt.Println("可用区: ", azName)

	//收集输出body
	data := make(map[string]string)
	data["当前可用区"] = azName

	//输出所有Env
	envs := os.Environ()
	for _, e := range envs {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		} else {
			data["Env-"+string(parts[0])] = string(parts[1])
		}
	}

	//结构化日志生成
	//debug := Debug{os.Getenv("CLOUD_ENV"), os.Getenv("flag"), azName}
	//a, _ := json.Marshal(debug)
	//os.Stdout.Write(a)

	//输出所有Header
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			data["Header-"+k] = v[0]
		}
	}

	//输出所有Form
	r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			data["Form-"+k] = v[0]
		}
	}

	//以下为设置返回（勿动）
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}

	url := strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
	data["url"] = url

	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}

	w.Write(msg)
}

// HeaderHandler 输出所有header
func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	//Count("/api/header")
	data := make(map[string]string)
	//输出所有Header
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			data[k] = v[0]
		}
	}

	//以下为设置返回（勿动）
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Write(msg)
}

// EnvHandler 输出所有env
func EnvHandler(w http.ResponseWriter, r *http.Request) {
	//Count("/api/env")
	data := make(map[string]string)

	//输出所有Env
	envs := os.Environ()
	for _, e := range envs {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		} else {
			data[string(parts[0])] = string(parts[1])
		}
	}

	//以下为设置返回（勿动）
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Write(msg)
}

// AzHandler 输出可用区
func AzHandler(w http.ResponseWriter, r *http.Request) {
	//Count("/api/az")
	data := make(map[string]string)

	var azName = getAzName(w)
	data["当前可用区"] = azName

	//以下为设置返回（勿动）
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Write(msg)
}

// OpenAPI 检查openapi https
func OpenAPI(w http.ResponseWriter, r *http.Request) {
	//Count("/api/https")
	appid := r.Header.Get("X-Tt-Appid")
	fmt.Println("调用openapi appid：", appid)

	url := "https://developer.toutiao.com/api/apps/qrcode"
	method := "POST"

	payload := strings.NewReader(`{
		"appname": "douyin"
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-DYC-FROM-APP", appid)
	req.Header.Add("X-DYC-SERVICE", r.Header.Get("X-Tt-Serviceid"))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	data := make(map[string]string)
	data["请求结果(20位)"] = string(body)[0:20]
	data["结论判断"] = "请求结果为PNG开头乱码即为【通过】，如果有token相关报错视为【不通过】"
	//输出所有Header
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			data[k] = v[0]
		}
	}

	//以下为设置返回（勿动）
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Write(msg)
}

// Http 检查openapi http
func Http(w http.ResponseWriter, r *http.Request) {
	//Count("/api/http")
	appid := r.Header.Get("X-Tt-Appid")
	fmt.Println("调用openapi appid：", appid)

	url := "http://developer.toutiao.com/api/apps/qrcode"
	method := "POST"

	payload := strings.NewReader(`{
		"appname": "douyin"
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-DYC-FROM-APP", appid)
	req.Header.Add("X-DYC-SERVICE", r.Header.Get("X-Tt-Serviceid"))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	data := make(map[string]string)
	data["请求结果(20位)"] = string(body)[0:20]
	data["结论判断"] = "请求结果为PNG开头乱码即为【通过】，如果有token相关报错视为【不通过】"
	//输出所有Header
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			data[k] = v[0]
		}
	}

	//以下为设置返回（勿动）
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Write(msg)
}

// PacketHandlerRequest 上行包大小
func PacketHandlerRequest(w http.ResponseWriter, r *http.Request) {
	//Count("/api/packetrequest")
	data := make(map[string]string)

	data["上行包大小"] = strconv.FormatInt(r.ContentLength, 10)
	data["上行上限为"] = "1M = 1048576字节"

	//以下为设置返回（勿动）
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Write(msg)
}

// PacketHandlerResponse 下行包大小
func PacketHandlerResponse(w http.ResponseWriter, r *http.Request) { //104848返回包大小 1048576=1M
	//Count("/api/packetresponse")
	base := 104848
	size := r.FormValue("size")
	if size == "" {
		size = "0"
	}
	sizeint, _ := strconv.Atoi(size)
	base += sizeint

	data := make(map[string]string)
	for i := 0; i < base; i++ {
		data[""] += "1111111111"
	}

	data["包大小约为"] = strconv.FormatInt(int64(base*10+90), 10) + "（精确包大小请在日志查看，1M == 1048576字节）"

	//以下为设置返回（勿动）
	msg, err := json.Marshal(data)
	fmt.Println("返回包大小", len(msg))
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	fmt.Fprintf(w, "%s", msg)
}

// Err 返回指定错误
func Err(w http.ResponseWriter, req *http.Request) {
	c := req.FormValue("code")
	code, _ := strconv.ParseInt(c, 10, 64)
	os.Stderr.WriteString("Msg to err\n")
	http.Error(w, "this is an err interface", int(code))
}

// Panic 返回panic
func Panic(w http.ResponseWriter, req *http.Request) {
	os.Stderr.WriteString("Msg to Panic\n")
	panic(req)
}

// OutLog 结构化日志
func OutLog(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	age, _ := strconv.ParseInt(req.FormValue("age"), 10, 64)
	salary, _ := strconv.ParseFloat(req.FormValue("salary"), 10)
	logID := req.FormValue("logID")
	height, _ := strconv.ParseFloat(req.FormValue("height"), 10)

	log.SetLevel(log.InfoLevel) // 设置输出警告级别
	if req.FormValue("format") == "text" {
		log.SetFormatter(&log.TextFormatter{}) // 设置 fromat text
	} else {
		log.SetFormatter(&log.JSONFormatter{}) // 设置 format json
	}

	if req.FormValue("format") == "text" {
		if req.FormValue("type") == "stdout" {
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(os.Stderr)
		}
		contextLogger := log.WithFields(log.Fields{
			"name":   name,
			"age":    age,
			"salary": salary,
		})

		contextLogger.Warn("这是一个warn级别日志")
		contextLogger.Info("这是一个info级别日志")
	} else {
		jLog := &jsonLog{
			LogID:  logID,
			Age:    age,
			Salary: salary,
			Name:   name,
			Height: height,
		}
		bytes, _ := json.Marshal(jLog)
		if req.FormValue("type") == "stdout" {
			os.Stdout.Write(bytes)
			fmt.Println()
		} else {
			os.Stderr.Write(bytes)
			fmt.Println()
		}
	}
}

// CloudId 返回info
func CloudId(w http.ResponseWriter, req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	fmt.Fprintf(w, "%s", s)
}
