package service

import (
	"douyincloud-gin-demo/db/mysql"
	"encoding/json"
	"fmt"
	"testing"
)

func TestUpdate(t *testing.T) {
	mysql.InitMysql()
	jsonString := "{\n        \"questionaireId\": \"1686139095162\",\n        \"title\": \"相似度分析测试中心\",\n        \"naireType\": 0,\n        \"iconUrl\": \"aiTools/icon1686139028517.png\",\n        \"iconTitle\": \"立即测试\",\n        \"homepageUrl\": \"aiTools/qusHomeBg1686139065073.png\",\n        \"ansertSheetUrl\": \"aiTools/answerSheetBg1686139072721.png\",\n        \"resultSheetUrl\": \"aiTools/resPageBg1686139080723.png\",\n        \"questions\": [\n                {\n                        \"questionaireId\": \"1686139095162\",\n                        \"questionId\": \"1686139095162_0\",\n                        \"content\": \"A\",\n                        \"answers\": [\n                                {\n                                        \"questionId\": \"1686139095162_0\",\n                                        \"answerId\": \"1686139095162_0_0\",\n                                        \"content\": \"A\"\n                                },\n                                {\n                                        \"questionId\": \"1686139095162_0\",\n                                        \"answerId\": \"1686139095162_0_1\",\n                                        \"content\": \"C\"\n                                }\n                        ],\n                        \"ownerAnswerId\": \"1686139095162_0_0\"\n                }\n        ]\n}"
	naire := &CreateQuestionnaireReq{}
	json.Unmarshal([]byte(jsonString), naire)

	TestUpdateFUnc(naire)

	fmt.Println(111)
}

func Test1(t *testing.T) {
	msg := fmt.Sprintf("[GetQuestionnaireInfo] begin naireId=%v", 11)
	fmt.Println(msg)
}
