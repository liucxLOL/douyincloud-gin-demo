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
