package service

import (
	"douyincloud-gin-demo/db/mysql"
	"encoding/json"
	"fmt"
	"testing"
)

func TestUpdate(t *testing.T) {
	mysql.InitMysql()
	jsonString := "{\n        \"questionaireId\": \"1686156464458\",\n        \"title\": \"新增问卷3\",\n        \"naireType\": 0,\n        \"iconUrl\": \"aiTools/icon1686156431069.png\",\n        \"iconTitle\": \"立即测试\",\n        \"homepageUrl\": \"\",\n        \"ansertSheetUrl\": \"\",\n        \"resultSheetUrl\": \"\",\n        \"questions\": [\n                {\n                        \"questionaireId\": \"1686156464458\",\n                        \"questionId\": \"1686156464458_0\",\n                        \"content\": \"问卷1\",\n                        \"answers\": [\n                                {\n                                        \"questionId\": \"1686156464458_0\",\n                                        \"answerId\": \"1686156464458_0_0\",\n                                        \"content\": \"a\"\n                                }\n                        ],\n                        \"ownerAnswerId\": \"1686156464458_0_0\"\n                },\n                {\n                        \"questionaireId\": \"1686156464458\",\n                        \"questionId\": \"1686156464458_1\",\n                        \"content\": \"问卷2\",\n                        \"answers\": [\n                                {\n                                        \"questionId\": \"1686156464458_1\",\n                                        \"answerId\": \"1686156464458_1_0\",\n                                        \"content\": \"2\"\n                                },\n                                {\n                                        \"questionId\": \"1686156464458_1\",\n                                        \"answerId\": \"1686156464458_1_1\",\n                                        \"content\": \"3\"\n                                }\n                        ],\n                        \"ownerAnswerId\": \"1686156464458_1_1\"\n                }\n        ]\n}"
	naire := &CreateQuestionnaireReq{}
	json.Unmarshal([]byte(jsonString), naire)

	TestUpdateFUnc(naire)

	fmt.Println(111)
}

func Test1(t *testing.T) {
	iamgePic := ""
	TestAiGetPic(1, iamgePic)
}
