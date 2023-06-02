package model

import (
	"douyincloud-gin-demo/db/mysql"
	"fmt"

	"gorm.io/gorm"
)

// 答案表
type Answer struct {
	Id         uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:自增id" json:"id"`
	AnswerId   string `gorm:"column:answer_id;type:varchar(64);comment:答案id;NOT NULL" json:"answer_id"`
	QuestionId string `gorm:"column:question_id;type:varchar(64);comment:问题id;NOT NULL" json:"question_id"`
	Content    string `gorm:"column:content;type:varchar(255);comment:结果内容;NOT NULL" json:"content"`
}

type AnswerDto struct {
	QuestionId string `json:"questionId"` // 问题id
	AnswerId   string `json:"answerId"`   //答案id
	Content    string `json:"content"`    //答案文本
}

func (a *Answer) TransAnswer2Dto() *AnswerDto {
	return &AnswerDto{
		AnswerId: a.AnswerId,
		Content:  a.Content,
	}
}

func (d *AnswerDto) TransAnswerDto2Answer() *Answer {
	return &Answer{
		QuestionId: d.QuestionId,
		AnswerId:   d.AnswerId,
		Content:    d.Content,
	}
}

func (m *Answer) TableName() string {
	return "answer"
}

const answerTableName = "answer"

// Get dbInstance
func GetMysql() *gorm.DB {
	return mysql.DbInstance
}

func SelectBy(id string) (*Answer, error) {
	db := GetMysql()
	var err error
	var model Answer
	err = db.Debug().Table(answerTableName).
		Where("id = ?", id).Scan(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func SelectAnswersByQuestionId(questionId string) ([]*Answer, error) {
	db := GetMysql()
	var err error
	var models []*Answer
	err = db.Debug().Table(answerTableName).
		Where("question_id = ?", questionId).Scan(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func InsertAnswer(model *Answer) error {
	db := GetMysql()
	var err error
	err = db.Debug().Table(answerTableName).
		Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateAnswer(model *AnswerDto) error {
	db := GetMysql()
	var err error
	err = db.Debug().Table(answerTableName).
		Where("answer_id = ?", model.AnswerId).Scan(&model).Error
	if err != nil {
		fmt.Sprintf("no log answerdto=%v", model)
		return err
	}

	err = db.Debug().Table(answerTableName).
		Where("answer_id = ?", model.TransAnswerDto2Answer()).Updates(&model).Error
	if err != nil {
		fmt.Sprintf("update answer faild answerdto=%v", model)
		return err
	}
	return nil
}