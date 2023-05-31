package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var ctx = context.Background()

// PingHandler 火山校验
func PingHandler(w http.ResponseWriter, r *http.Request) {
	Count("1")
	fmt.Fprintf(w, "pong!\n")
}

// HelpHandler 帮助
func HelpHandler(w http.ResponseWriter, r *http.Request) {
	//Count("/help")
	data := make(map[string]string)
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, "service", nil, 0)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		os.Exit(1)
	}

	var funcs []*ast.FuncDecl
	for _, pack := range packs {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					funcs = append(funcs, fn)
				}
			}
		}
	}
	for _, tmp := range funcs {
		data[tmp.Name.Name] = tmp.Doc.Text()
	}
	data["DEMO功能详情见文档"] = ""
	//以下为设置返回（勿动）
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Write(msg)
}

func TestHttp(w http.ResponseWriter, r *http.Request) {
	url := "http://cloud-database-api.dycloud.run/api/cloud_database/exec_cloud_database_cmd"
	method := "POST"

	payload := strings.NewReader(`{
    "collection_name": "todos",
    "query_type": "WHERE",
    "offset": 0,
    "limit": 100,
    "action": "database.getDocument"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

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
	return
}

func Count(count string) { //上报部署数据
	//res, _ := http.Get("https://cloudapi.bytedance.net/faas/services/tt7g27/invoke/data_TestDemo?count=" + count)
	//fmt.Println("上报数据：\n", res)
	//var (
	//	client     *mongo.Client
	//	err        error
	//	db         *mongo.Database
	//	collection *mongo.Collection
	//)
	//type Count struct {
	//	StartTime int64 `bson:"startTime"` //开始时间
	//	EndTime   int64 `bson:"endTime"`   //结束时间
	//}
	////1.建立连接
	//if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:Aa12345!@mongoreplicacc99660f06f10.mongodb.volces.com:3717,mongoreplicacc99660f06f11.mongodb.volces.com:3717/?authSource=admin&replicaSet=rs-mongo-replica-cc99660f06f1&retryWrites=true").SetConnectTimeout(5*time.Second)); err != nil {
	//	fmt.Print(err)
	//	fmt.Print(client)
	//	return
	//}
	////2.选择数据库 my_db
	//db = client.Database("qa_db")
	//fmt.Print("dbname: ", db.Name())
	//
	////3.选择表 my_collection
	//collection = db.Collection("testdemo")
	//collection = collection
	//fmt.Print("collection: ", collection.Name())
	//
	////
	//collection.findb(context.TODO(), bson.D{{"_id", id}}, opts)
}
