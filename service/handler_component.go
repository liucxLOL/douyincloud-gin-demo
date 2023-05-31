package service

import (
	"douyincloud-gin-demo/db/mongodb"
	"douyincloud-gin-demo/db/mysql"
	"douyincloud-gin-demo/db/redis"
	"douyincloud-gin-demo/db/tos"
	"encoding/json"
	"fmt"

	"net/http"
)

func MysqlSelect(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	model, err := mysql.Select(id)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", model, err)
}

func MysqlSelectList(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	models, err := mysql.SelectList(name)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", models, err)
}

func MysqlCreate(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	model, err := mysql.Create(name)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", model, err)
}

func MysqlCreateLockTable(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	model, err := mysql.CreateLockTable(name)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", model, err)
}

func MysqlUpdate(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	model, err := mysql.Update(id)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", model, err)
}

func MysqlUpdateCounts(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	count := req.FormValue("count")
	err := mysql.UpdateCounts(name, count)
	fmt.Fprintf(w, "err: %v\n", err)
}

func MysqlDelete(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	err := mysql.Delete(id)
	fmt.Fprintf(w, "err: %v\n", err)
}

func MysqlDeleteRollback(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	err := mysql.DeleteRollback(id)
	fmt.Fprintf(w, "err: %v\n", err)
}

func RedisSet(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")
	value := req.FormValue("value")
	expireTime := req.FormValue("expireTime")
	res, err := redis.Set(key, value, expireTime)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", res, err)
}

func RedisGet(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")
	res, err := redis.Get(key)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", res, err)
}

func RedisDel(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")
	res, err := redis.Del(key)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", res, err)
}

func MongoInsert(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	count := req.FormValue("count")
	res, err := mongodb.InsertOne(name, count)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", res, err)
}

func MongoFind(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	count := req.FormValue("count")
	res, err := mongodb.FindOne(name, count)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", res, err)
}

func MongoDelete(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	count := req.FormValue("count")
	res, err := mongodb.DeleteOne(name, count)
	fmt.Fprintf(w, "Response: %+v, err: %v\n", res, err)
}

func TosPutObject(w http.ResponseWriter, r *http.Request) {
	// 获取报文的长度
	len := r.ContentLength
	// 定义一个byte切片
	data := make([]byte, len)
	// 读取请求体
	r.Body.Read(data)
	body := make(map[string]string)
	json.Unmarshal(data, &body)
	res, err := tos.PutObject(body["endpoint"], body["accessKey"], body["secretKey"], body["bucketName"], body["objectKey"])
	fmt.Fprintf(w, "Response: %+v, err: %v\n", res, err)

}
func TosGetObject(w http.ResponseWriter, r *http.Request) {
	// 获取报文的长度
	len := r.ContentLength
	// 定义一个byte切片
	data := make([]byte, len)
	// 读取请求体
	r.Body.Read(data)
	body := make(map[string]string)
	json.Unmarshal(data, &body)

	res, err := tos.GetObject(body["endpoint"], body["accessKey"], body["secretKey"], body["bucketName"], body["objectKey"])
	fmt.Fprintf(w, "Response: %+v, err: %v\n", res, err)
}

func TosUploadPart(w http.ResponseWriter, r *http.Request) {
	// 获取报文的长度
	len := r.ContentLength
	// 定义一个byte切片
	data := make([]byte, len)
	// 读取请求体
	r.Body.Read(data)
	body := make(map[string]string)
	json.Unmarshal(data, &body)

	res, err := tos.UploadPart(body["endpoint"], body["accessKey"], body["secretKey"], body["bucketName"], body["objectKey"])
	fmt.Fprintf(w, "Response: %+v, err: %v\n", res, err)

}
