package service

import (
	"context"
	"douyincloud-gin-demo/db/mysql/model"
	"douyincloud-gin-demo/service/handle_volc"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type VolcAiReq struct {
	ImageBase64 string `json:"image_base_64"`
	Type        int    `json:"type"` //1.漫画风，2老照片修复
}

type VolcAiRespon struct {
	ImageBase64 string `json:"image_base_64"`
}

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
	log.Info("[SelectQuestionnaireList] owner =" + owner + ",openId=" + openID)
	ctx := context.Background()
	naires := []*model.Questionnaire{}
	err := errors.New("systemError")
	//如果为true
	if owner == "true" {
		naires, err = model.SelectQuestionnaireByOpenId(openID)
		if err != nil || len(naires) == 0 {
			log.Error(fmt.Sprintf("True getnaire faild err=%v,naires=%v", err, naires))
			FillResponse(ctx, w, 0, nil)
			return
		}
	} else {
		naires, err = model.SelectQuestionnaires()
		if err != nil || len(naires) == 0 {
			log.Error(fmt.Sprintf("False getnaire faild err=%v,naires=%v", err, naires))
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
	log.Info(fmt.Sprintf("[GetQuestionnaireInfo] begin naireId=%v", naireId))
	naireInfo, err := model.SelectQuestionnaireById(naireId)
	if err != nil {
		log.Error(fmt.Sprintf("[GetQuestionnaireInfo] faild naireId=%v", naireId))
		FillResponse(ctx, w, 1, nil)
		return
	}
	if naireInfo == nil {
		log.Warning(fmt.Sprintf("[GetQuestionnaireInfo] none naireId=%v", naireId))
		FillResponse(ctx, w, 0, nil)
		return
	}

	naireId = naireInfo.QuestionaireId

	questions, err := model.SelectQuestionByQuestionNaireId(naireId)
	if err != nil {
		log.Error(fmt.Sprintf("[SelectQuestionByQuestionNaireId] faild naireId=%v", naireId))
		FillResponse(ctx, w, 1, nil)
		return
	}

	if len(questions) == 0 {
		log.Error(fmt.Sprintf("[SelectQuestionByQuestionNaireId] question is none  naireId=%v", naireId))
		FillResponse(ctx, w, 1001, "question is nil")
		return
	}

	questionResps := []Question{}

	for _, question := range questions {
		answerResps := []Answer{}
		questionId := question.QuestionId
		answers, err := model.SelectAnswersByQuestionId(questionId)
		if err != nil {
			log.Error(fmt.Sprintf("[SelectAnswersByQuestionId] faild naireId=%v", naireId))
			FillResponse(ctx, w, 1, nil)
			return
		}

		if len(answers) == 0 {
			log.Error(fmt.Sprintf("[SelectAnswersByQuestionId] answers is none  naireId=%v", naireId))
			FillResponse(ctx, w, 1, "answer is nil")
			return
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

func UpdateQuestionnaireInfo(w http.ResponseWriter, req *http.Request) {
	log.Info("UpdateQuestionnaireInfo begin")
	ctx := context.Background()
	naireReq := &CreateQuestionnaireReq{}
	err := json.NewDecoder(req.Body).Decode(naireReq)
	b, _ := json.Marshal(naireReq)
	log.Info(fmt.Sprintf("[UpdateQuestionnaireInfo] req info=%v", string(b)))
	openID := req.Header.Get("X-TT-OPENID")
	if err != nil || openID == "" {
		log.Error(fmt.Sprintf("[UpdateQuestionnaireInfo] trans req 2 model faild err=%v", err))
		FillResponse(ctx, w, 1, nil)
		return
	}

	if len(naireReq.Questions) == 0 {
		FillResponse(ctx, w, 1001, "quesiton number is 0")
		return
	}

	useQuestionIds := []string{}
	useAnswerIds := []string{}

	//增加或者更新
	for _, question := range naireReq.Questions {
		useAnswerIds = []string{}
		//保存Answers
		if len(question.Answers) == 0 {
			FillResponse(ctx, w, 1002, "answer number is 0")
			return
		}
		for _, answer := range question.Answers {
			useAnswerIds = append(useAnswerIds, answer.AnswerId)
			answerModel := &model.Answer{
				AnswerId:   answer.AnswerId,
				QuestionId: answer.QuestionId,
				Content:    answer.Content,
			}
			err := model.UpdateAnswer(answerModel)
			if err != nil {
				log.Error(fmt.Sprintf("update answer faild err=%v", err))
				FillResponse(ctx, w, 1, nil)
				return
			}

		}

		//保存Questions
		questionModel := &model.Question{
			QuestionId:     question.QuestionId,
			Content:        question.Content,
			AnswerId:       question.OwnerAnswerId,
			QuestionaireId: question.QuestionaireId,
		}
		err := model.UpdateQuestion(questionModel)
		if err != nil {
			log.Error(fmt.Sprintf("update  question faild err=%v", err))
			FillResponse(ctx, w, 1, nil)
			return
		}

		useQuestionIds = append(useQuestionIds, question.QuestionId)
		log.Info(fmt.Sprintf("del questionId=%v,answerIds=%v", question.QuestionId, useAnswerIds))
		//删除多余的answer
		model.DelAnswerNotInUse(useAnswerIds)
	}

	log.Info(fmt.Sprintf("del questionIds=%v", useQuestionIds))
	//删除多余的question
	model.DelQuestonNotInUse(useQuestionIds)

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

	isContinue := SetQuestionnaires(ctx, naireReq, true)

	if !isContinue {
		FillResponse(ctx, w, 2001, nil)
		return
	}

	err = model.UpdateQuestionnaire(naireModel)
	if err != nil {
		log.Error(fmt.Sprintf("update  questionnaire faild err=%v", err))
		FillResponse(ctx, w, 1, nil)
		return
	}

	FillResponse(ctx, w, 0, nil)

}

func TestUpdateFUnc(naireReq *CreateQuestionnaireReq) {
	log.Info("UpdateQuestionnaireInfo begin")
	ctx := context.Background()

	useQuestionIds := []string{}
	useAnswerIds := []string{}

	//增加或者更新
	for _, question := range naireReq.Questions {
		//保存Answers
		for _, answer := range question.Answers {
			useAnswerIds = append(useAnswerIds, answer.AnswerId)
			answerModel := &model.Answer{
				AnswerId:   answer.AnswerId,
				QuestionId: answer.QuestionId,
				Content:    answer.Content,
			}
			err := model.UpdateAnswer(answerModel)
			if err != nil {
				log.Error(fmt.Sprintf("update answer faild err=%v", err))

				return
			}

		}

		//保存Questions
		questionModel := &model.Question{
			QuestionId:     question.QuestionId,
			Content:        question.Content,
			AnswerId:       question.OwnerAnswerId,
			QuestionaireId: question.QuestionaireId,
		}
		err := model.UpdateQuestion(questionModel)
		if err != nil {
			log.Error(fmt.Sprintf("update  question faild err=%v", err))

			return
		}

		useQuestionIds = append(useQuestionIds, question.QuestionId)
		log.Info(fmt.Sprintf("del questionId=%v,answerIds=%v", question.QuestionId, useAnswerIds))
	}

	log.Info(fmt.Sprintf("del questionIds=%v", useQuestionIds))
	//删除多余的question
	model.DelQuestonNotInUse(useQuestionIds)
	//删除多余的answer
	model.DelAnswerNotInUse(useAnswerIds)

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
	}

	isContinue := SetQuestionnaires(ctx, naireReq, true)

	if !isContinue {

		return
	}

	err := model.UpdateQuestionnaire(naireModel)
	if err != nil {
		log.Error(fmt.Sprintf("update  questionnaire faild err=%v", err))

		return
	}

}

func CreateQuestionnaireInfo(w http.ResponseWriter, req *http.Request) {
	log.Info("CreateQuestionnaireInfo begin")
	ctx := context.Background()
	naireReq := &CreateQuestionnaireReq{}
	err := json.NewDecoder(req.Body).Decode(naireReq)
	openID := req.Header.Get("X-TT-OPENID")
	if err != nil || openID == "" {
		log.Error(fmt.Sprintf("[CreateQuestionnaireInfo] trans req 2 model faild err=%v", err))
		FillResponse(ctx, w, 1, nil)
		return
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
				log.Error(fmt.Sprintf("insert into answer faild err=%v", err))
				FillResponse(ctx, w, 1, nil)
				return
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
			log.Error(fmt.Sprintf("insert into question faild err=%v", err))
			FillResponse(ctx, w, 1, nil)
			return
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

	isContinue := SetQuestionnaires(ctx, naireReq, false)

	if !isContinue {
		FillResponse(ctx, w, 2001, nil)
		return
	}

	log.Info(ctx, fmt.Sprintf("naireReqId=%v, modelId=%v", naireReq.QuestionaireId, naireModel.QuestionaireId))
	err = model.InsertQuestionnaire(naireModel)
	if err != nil {
		log.Error(fmt.Sprintf("insert into questionnaire faild err=%v", err))
		FillResponse(ctx, w, 1, nil)
		return
	}

	FillResponse(ctx, w, 0, nil)
}

func VolcAIGetPic(w http.ResponseWriter, req *http.Request) {
	log.Info("VolcAIpicBegin")
	ctx := context.Background()
	volcAiReq := &VolcAiReq{}
	err := json.NewDecoder(req.Body).Decode(volcAiReq)
	retImage := ""
	if err != nil {
		log.Error(fmt.Sprintf("[VolcAIGetPic] trans req 2 model faild err=%v", err))
		FillResponse(ctx, w, 1, nil)
		return
	}
	if volcAiReq.Type == 1 {
		//漫画风
		retImage = handle_volc.GetAIPhoto(volcAiReq.ImageBase64, 1)

	} else if volcAiReq.Type == 2 {
		//老照片修复
		retImage = handle_volc.GetAIPhoto(volcAiReq.ImageBase64, 2)
	}

	if retImage == "" {
		log.Info("VolcAIGetPic end retImage faild")
		FillResponse(ctx, w, 1, "")
		return
	}

	FillResponse(ctx, w, 0, retImage)
}

func SetQuestionnaires(ctx context.Context, req *CreateQuestionnaireReq, isUpdate bool) bool {
	urls := []string{}
	if !isUpdate {

		if req.HomepageUrl != "" {
			urls = append(urls, req.HomepageUrl)
		}

		if req.IconUrl != "" {
			urls = append(urls, req.IconUrl)
		}

		if req.ResultSheetUrl != "" {
			urls = append(urls, req.ResultSheetUrl)
		}

		if req.AnsertSheetUrl != "" {
			urls = append(urls, req.AnsertSheetUrl)
		}

		if len(urls) == 0 {
			return false
		}
	} else {
		naireId := req.QuestionaireId
		naire, err := model.SelectQuestionnaireById(naireId)
		if err != nil || naire == nil {
			log.Error(fmt.Sprintf("update naire faild naireId=%v", naireId))
			return false
		}

		if naire.ResultSheetUrl != req.ResultSheetUrl {
			urls = append(urls, req.ResultSheetUrl)
		}

		if naire.AnsertSheetUrl != req.AnsertSheetUrl {
			urls = append(urls, req.AnsertSheetUrl)
		}

		if naire.IconUrl != req.IconUrl {
			urls = append(urls, req.IconUrl)
		}

		if naire.HomepageUrl != req.HomepageUrl {
			urls = append(urls, req.HomepageUrl)
		}
	}
	return handle_volc.SetPicPublic(ctx, urls)
}
func GetModelFromReq(ctx context.Context, req *http.Request, model interface{}) interface{} {
	err := json.NewDecoder(req.Body).Decode(model)
	if err != nil {
		log.Error(fmt.Sprintf("trans req 2 model faild err=%v", err))
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
