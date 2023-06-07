package service

import (
	"douyincloud-gin-demo/db/mysql"
	"encoding/json"
	"fmt"
	"testing"
)

func TestUpdate(t *testing.T) {
	mysql.InitMysql()
	jsonString := "{\n\t\t\"questionaireId\": \"1686139095162\",\n\t\t\"title\": \"相似度分析测试中心\",\n\t\t\"naireType\": 0,\n\t\t\"iconUrl\": \"aiTools/icon1686139028517.png\",\n\t\t\"iconTitle\": \"立即测试\",\n\t\t\"homepageUrl\": \"aiTools/qusHomeBg1686139065073.png\",\n\t\t\"ansertSheetUrl\": \"aiTools/answerSheetBg1686139072721.png\",\n\t\t\"resultSheetUrl\": \"aiTools/resPageBg1686139080723.png\",\n\t\t\"questions\": [\n\t{\n\t\"questionaireId\": \"1686139095162\",\n\t\"questionId\": \"1686139095162_0\",\n\t\"content\": \"A\",\n\t\"answers\": [\n\t{\n\t\"questionId\": \"1686139095162_0\",\n\t\"answerId\": \"1686139095162_0_0\",\n\t\"content\": \"A\"\n\t},\n\t{\n\t\"questionId\": \"1686139095162_0\",\n\t\"answerId\": \"1686139095162_0_1\",\n\t\"content\": \"C\"\n\t}\n\t],\n\t\"ownerAnswerId\": \"1686139095162_0_0\"\n\t}\n]\n}"
	naire := &CreateQuestionnaireReq{}
	json.Unmarshal([]byte(jsonString), naire)

	TestUpdateFUnc(naire)

	fmt.Println(111)
}
