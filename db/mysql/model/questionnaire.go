package model

import "fmt"

// 问卷信息表
type Questionnaire struct {
	Id             uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:自增id" json:"id"`
	QuestionaireId string `gorm:"column:questionaire_id;type:varchar(64);comment:问卷id;NOT NULL" json:"questionaire_id"`
	Title          string `gorm:"column:title;type:varchar(32);comment:标题;NOT NULL" json:"title"`
	Type           int    `gorm:"column:type;type:int(11);comment:0-相似度，1-ai作画;NOT NULL" json:"type"`
	IconUrl        string `gorm:"column:icon_url;type:varchar(64);comment:图标url;NOT NULL" json:"icon_url"`
	IconTitle      string `gorm:"column:icon_title;type:varchar(64);comment:开始按钮的标题;NOT NULL" json:"icon_title"`
	HomepageUrl    string `gorm:"column:homepage_url;type:varchar(64);comment:问卷主页背景图;NOT NULL" json:"homepage_url"`
	AnsertSheetUrl string `gorm:"column:ansert_sheet_url;type:varchar(64);comment:答题页背景图;NOT NULL" json:"ansert_sheet_url"`
	ResultSheetUrl string `gorm:"column:result_sheet_url;type:varchar(64);comment:结果页背景图;NOT NULL" json:"result_sheet_url"`
	CreatorOpenId  string `gorm:"column:creator_open_id;type:varchar(64);comment:问卷创作者;NOT NULL" json:"creator_open_id"`
	QuestionList   string `gorm:"column:question_list;type:varchar(64);comment:包含的问题列表;NOT NULL" json:"question_list"`
}

func (m *Questionnaire) TableName() string {
	return "questionnaire"
}

const QuestionnaireTableName = "questionnaire"

func SelectQuestionnaireByOpenId(openId string) ([]*Questionnaire, error) {
	db := GetMysql()
	var err error
	var models []*Questionnaire
	err = db.Debug().Table(QuestionnaireTableName).
		Where("creator_open_id = ?", openId).Scan(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func SelectQuestionnaireById(QuestionaireId string) (*Questionnaire, error) {
	db := GetMysql()
	var err error
	var model Questionnaire
	err = db.Debug().Table(QuestionnaireTableName).
		Where("questionaire_id = ?", QuestionaireId).Scan(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func SelectQuestionnaires() ([]*Questionnaire, error) {
	db := GetMysql()
	var err error
	var models []*Questionnaire
	err = db.Debug().Table(QuestionnaireTableName).Scan(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func InsertQuestionnaire(model *Questionnaire) error {
	db := GetMysql()
	var err error
	err = db.Debug().Table(QuestionnaireTableName).
		Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateQuestionnaire(model *Questionnaire) error {
	db := GetMysql()
	var err error
	err = db.Debug().Table(QuestionnaireTableName).
		Where("questionaire_id = ?", model.QuestionaireId).Scan(&model).Error
	if err != nil {
		fmt.Sprintf("no log questionnaire=%v", model)
		return err
	}

	err = db.Debug().Table(QuestionnaireTableName).
		Where("questionaire_id = ?", model.QuestionaireId).Updates(&model).Error
	if err != nil {
		fmt.Sprintf("update questionnaire faild questionnaire=%v", model)
		return err
	}
	return nil
}
