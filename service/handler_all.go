package service

import (
	"context"
	"douyincloud-gin-demo/db/mysql/model"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Questionnaire struct {
	QuestionaireId string `json:"questionaireId"` //问卷id
	Title          string `json:"title"`          //问卷标题
	IconUrl        string `json:"iconUrl"`        //图标
	NaireType      int32  `json:"naireType"`      //问卷类型
}

type QuestionnaireInfo struct {
	QuestionaireId string     `json:"questionaireId"` // '问卷id'
	Title          string     `json:"title"`          //  '标题',
	IconTitle      string     `json:"iconTitle"`      // '按钮文案'
	NaireType      int32      `json:"naireType"`      //'0-相似度，1-ai作画',
	IconUrl        string     `json:"iconUrl"`        //'图标url',
	HomepageUrl    string     `json:"homepageUrl"`    //'问卷主页背景图'
	AnsertSheetUrl string     `json:"ansertSheetUrl"` //'答题页背景图',
	ResultSheetUrl string     `json:"resultSheetUrl"` //'结果页背景图',
	Questions      []Question `json:"questions"`      //问题list
}

type Question struct {
	QuestionaireId string   `json:"questionaireId"` // '问卷id'
	QuestionId     string   `json:"questionId"`     // 问题id
	Content        string   `json:"content"`        // 问题内容
	OwnerAnswerId  string   `json:"ownerAnswerId"`  //达人的答案
	Answers        []Answer `json:"answers"`        //答案
	UserAnswerId   string   `json:"userAnswerId"`   //用户选择的答案
}

type Answer struct {
	QuestionId string `json:"questionId"` // 问题id
	AnswerId   string `json:"answerId"`   //答案id
	Content    string `json:"content"`    //答案文本
}

type CreateQuestionnaireReq struct {
	QuestionaireId string     `json:"questionaireId"` // '问卷id'
	Title          string     `json:"title"`          //  '标题'
	NaireType      int32      `json:"naireType"`      //'0-相似度，1-ai作画',
	IconUrl        string     `json:"iconUrl"`        //'图标url',
	IconTitle      string     `json:"iconTitle"`
	HomepageUrl    string     `json:"homepageUrl"`    //'问卷主页背景图'
	AnsertSheetUrl string     `json:"ansertSheetUrl"` //'答题页背景图',
	ResultSheetUrl string     `json:"resultSheetUrl"` //'结果页背景图',
	Questions      []Question `json:"questions"`      //问题list
}

func SelectQuestionnaireList(w http.ResponseWriter, req *http.Request) {
	owner := req.FormValue("owner")
	openID := req.Header.Get("X-TT-OPENID")
	ctx := context.Background()
	naires := []*model.Questionnaire{}
	//如果为true
	if owner == "true" {
		naires, err := model.SelectQuestionnaireByOpenId(openID)
		if err != nil || len(naires) == 0 {
			log.Error("getnaire faild err=%v,naires=%v", err, naires)
			FillResponse(ctx, w, 0, nil)
			return
		}
	} else {
		naires, err := model.SelectQuestionnaires()
		if err != nil || len(naires) == 0 {
			log.Error("getnaire faild err=%v,naires=%v", err, naires)
			FillResponse(ctx, w, 0, nil)
			return
		}
	}

	resp := []*Questionnaire{}

	for _, naire := range naires {
		naireInfo := &Questionnaire{
			QuestionaireId: naire.QuestionaireId,
			Title:          naire.Title,
			IconUrl:        naire.IconUrl,
			NaireType:      int32(naire.Type),
		}

		resp = append(resp, naireInfo)
	}

	FillResponse(ctx, w, 0, resp)
}

func GetQuestionnaireInfo(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	naireId := req.FormValue("questionaireId")
	log.Info("[GetQuestionnaireInfo] begin naireId=%v", naireId)
	naireInfo, err := model.SelectQuestionnaireById(naireId)
	if err != nil {
		log.Error("[GetQuestionnaireInfo] faild naireId=%v", naireId)
		FillResponse(ctx, w, 1, nil)
	}
	if naireInfo == nil {
		log.Warning("[GetQuestionnaireInfo] none naireId=%v", naireId)
		FillResponse(ctx, w, 0, nil)
	}

	naireId = naireInfo.QuestionaireId

	questions, err := model.SelectQuestionByQuestionNaireId(naireId)
	if err != nil {
		log.Error("[SelectQuestionByQuestionNaireId] faild naireId=%v", naireId)
		FillResponse(ctx, w, 1, nil)
	}

	if len(questions) == 0 {
		log.Error("[SelectQuestionByQuestionNaireId] question is none  naireId=%v", naireId)
		FillResponse(ctx, w, 1, nil)
	}

	questionResps := []Question{}

	for _, question := range questions {
		answerResps := []Answer{}
		questionId := question.QuestionId
		answers, err := model.SelectAnswersByQuestionId(questionId)
		if err != nil {
			log.Error("[SelectAnswersByQuestionId] faild naireId=%v", naireId)
			FillResponse(ctx, w, 1, nil)
		}

		if len(answers) == 0 {
			log.Error("[SelectAnswersByQuestionId] answers is none  naireId=%v", naireId)
			FillResponse(ctx, w, 1, nil)
		}

		for _, answer := range answers {
			answerResp := Answer{
				QuestionId: answer.QuestionId,
				AnswerId:   answer.AnswerId,
				Content:    answer.Content,
			}

			answerResps = append(answerResps, answerResp)
		}

		questionResp := Question{
			QuestionaireId: question.QuestionaireId,
			QuestionId:     question.QuestionId,
			Content:        question.Content,
			OwnerAnswerId:  question.AnswerId,
			Answers:        answerResps,
		}

		questionResps = append(questionResps, questionResp)
	}

	naireResp := QuestionnaireInfo{
		QuestionaireId: naireInfo.QuestionaireId,
		Questions:      questionResps,
		Title:          naireInfo.Title,
		IconUrl:        naireInfo.IconUrl,
		IconTitle:      naireInfo.IconTitle,
		NaireType:      int32(naireInfo.Type),
		HomepageUrl:    naireInfo.HomepageUrl,
		AnsertSheetUrl: naireInfo.AnsertSheetUrl,
		ResultSheetUrl: naireInfo.ResultSheetUrl,
	}

	FillResponse(ctx, w, 0, naireResp)

}

func CreateQuestionnaireInfo(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	naireReq := &CreateQuestionnaireReq{}
	err := json.NewDecoder(req.Body).Decode(naireReq)
	openID := req.Header.Get("X-TT-OPENID")
	if err != nil || openID == "" {
		log.Error("[CreateQuestionnaireInfo] trans req 2 model faild err=%v", err)
		FillResponse(ctx, w, 1, nil)
	}

	for _, question := range naireReq.Questions {
		//保存Answers
		for _, answer := range question.Answers {
			answerModel := &model.Answer{
				AnswerId:   answer.AnswerId,
				QuestionId: answer.QuestionId,
				Content:    answer.Content,
			}
			err := model.InsertAnswer(answerModel)
			if err != nil {
				log.Error("insert into answer faild err=%v", err)
				FillResponse(ctx, w, 1, nil)
			}

		}

		//保存Questions
		questionModel := &model.Question{
			QuestionId:     question.QuestionId,
			Content:        question.Content,
			AnswerId:       question.OwnerAnswerId,
			QuestionaireId: question.QuestionaireId,
		}
		err := model.InsertQuestion(questionModel)
		if err != nil {
			log.Error("insert into question faild err=%v", err)
			FillResponse(ctx, w, 1, nil)
		}
	}

	//保存Questionnaire
	naireModel := &model.Questionnaire{
		QuestionaireId: naireReq.QuestionaireId,
		Title:          naireReq.Title,
		Type:           int(naireReq.NaireType),
		IconUrl:        naireReq.IconUrl,
		IconTitle:      naireReq.IconTitle,
		HomepageUrl:    naireReq.HomepageUrl,
		AnsertSheetUrl: naireReq.AnsertSheetUrl,
		ResultSheetUrl: naireReq.ResultSheetUrl,
		CreatorOpenId:  openID,
	}
	err = model.InsertQuestionnaire(naireModel)
	if err != nil {
		log.Error("insert into questionnaire faild err=%v", err)
		FillResponse(ctx, w, 1, nil)
	}
}
func GetModelFromReq(ctx context.Context, req *http.Request, model interface{}) interface{} {
	err := json.NewDecoder(req.Body).Decode(model)
	if err != nil {
		log.Error("trans req 2 model faild err=%v", err)
		return nil
	}

	return model
}

func FillResponse(ctx context.Context, w http.ResponseWriter, statusCode int64, data interface{}) {
	w.WriteHeader(http.StatusOK)
	extra := make(map[string]interface{})
	extra["now"] = time.Now().UnixMilli()
	body := &ResponseWrapper{
		StatusCode: statusCode,
		Data:       data,
	}
	rawBody, e := json.Marshal(body)
	if e != nil {
		log.Error(ctx, "fail to marshal resp:%+v, e:%+v", body, e)
	}
	w.Write(rawBody)
}

type ResponseWrapper struct {
	StatusCode int64       `json:"code"`
	Data       interface{} `json:"data"`
}
