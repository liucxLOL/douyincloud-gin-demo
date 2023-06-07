package model

import (
	"douyincloud-gin-demo/db/mysql"
	"fmt"
	"testing"
)

func TestSelect(t *testing.T) {
	mysql.InitMysql()

	models, _ := SelectQuestionnaireById("1686139095162")
	fmt.Println(models)
	t.Log("success")
}

func TestUpdateAnswer(t *testing.T) {
	mysql.InitMysql()
	ans := &Answer{
		AnswerId:   "1686139095162_0_1",
		QuestionId: "liucx_1",
		Content:    "update_lcx_content",
	}
	UpdateAnswer(ans)

}
