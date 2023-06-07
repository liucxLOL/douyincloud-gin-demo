package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm/clause"
)

// 问题表
type Question struct {
	Id             uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:自增id" json:"id"`
	QuestionId     string `gorm:"column:question_id;type:varchar(64);comment:问题id;NOT NULL" json:"question_id"`
	Content        string `gorm:"column:content;type:varchar(255);comment:问题内容;NOT NULL" json:"content"`
	AnswerId       string `gorm:"column:answer_id;type:varchar(64);comment:正确答案id" json:"answer_id"`
	QuestionaireId string `gorm:"column:questionaire_id;type:varchar(64);comment:所属哪一个问卷;NOT NULL" json:"questionaire_id"`
}

type QuestionDto struct {
	QuestionaireId string   `json:"questionaireId"` // '问卷id'
	QuestionId     string   `json:"questionId"`     // 问题id
	Content        string   `json:"content"`        // 问题内容
	OwnerAnswerId  string   `json:"ownerAnswerId"`  //达人的答案
	Answers        []Answer `json:"answers"`        //答案
	UserAnswerId   string   `json:"userAnswerId"`   //用户选择的答案
}

func (m *Question) TableName() string {
	return "question"
}

const QuestionTableName = "question"

func (a *Question) TransQuestion2Dto() *QuestionDto {

	dto := &QuestionDto{
		QuestionaireId: a.QuestionaireId,
		QuestionId:     a.QuestionId,
		Content:        a.Content,
		OwnerAnswerId:  a.AnswerId,
	}
	return dto

}

func (d *QuestionDto) TransQuestionDto2Question() *Question {
	return &Question{
		QuestionId:     d.QuestionId,
		QuestionaireId: d.QuestionaireId,
		Content:        d.Content,
		AnswerId:       d.OwnerAnswerId,
	}
}

func SelectQuestionByQuestionNaireId(questionaireId string) ([]*Question, error) {
	db := GetMysql()
	var err error
	var models []*Question
	err = db.Debug().Table(QuestionTableName).
		Where("questionaire_id = ?", questionaireId).Scan(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func InsertQuestion(model *Question) error {
	db := GetMysql()
	var err error
	err = db.Debug().Table(QuestionTableName).
		Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateQuestion(model *Question) error {
	db := GetMysql()

	err := db.Table(QuestionTableName).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},                                                                // key colume
		DoUpdates: clause.AssignmentColumns([]string{"question_id", "content", "answer_id", "questionaire_id"}), // column needed to be updated
	}).Create(&model)

	if err != nil {
		fmt.Sprintf("update QuestionDto faild QuestionDto=%v", model)
		return errors.New("system error")
	}
	return nil
}
