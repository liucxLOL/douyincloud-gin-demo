package handle_volc

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/volcengine/volc-sdk-golang/service/visual"
	"github.com/volcengine/volc-sdk-golang/service/visual/model"
)

func GetAIPhotoStr() string {
	testAk := "AKLTMTdmNzQxZjczZjk0NGVkYTk3YzdiZmY5YWEyMTM1MGE"
	testSk := "T0dZek5EbGtaakkzTVdFNU5HTXdaRGc1TXpSbU56aG1aRGRoTURReU1UTQ=="
	visual.DefaultInstance.Client.SetAccessKey(testAk)
	visual.DefaultInstance.Client.SetSecretKey(testSk)
	//visual.DefaultInstance.SetRegion("region")
	//visual.DefaultInstance.SetHost("host")

	//请求入参
	reqBody := &model.T2ILDMRequest{
		ReqKey:    "t2i_ldm", // 固定值
		Text:      "家庭幸福是多方面的，其中一个重要方面是健康。保持苗条的身材不仅有助于身体健康，还能提高自信心，增强夫妻之间的吸引力，促进婚姻幸福。因此，我们应该注重饮食健康，多运动，保持好身材，从而为家庭幸福添彩",
		StyleTerm: "",
	}

	resp, status, err := visual.DefaultInstance.T2ILDM(reqBody)
	fmt.Println(status, err)
	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
	return ""
}

func GetAIPhoto(imageBase64 string, aiType int) (string, bool) {
	testAk := "AKLTMTdmNzQxZjczZjk0NGVkYTk3YzdiZmY5YWEyMTM1MGE"
	testSk := "T0dZek5EbGtaakkzTVdFNU5HTXdaRGc1TXpSbU56aG1aRGRoTURReU1UTQ=="

	visual.DefaultInstance.Client.SetAccessKey(testAk)
	visual.DefaultInstance.Client.SetSecretKey(testSk)

	form := url.Values{}
	form.Add("image_base64", imageBase64)

	if aiType == 1 {
		resp, status, err := visual.DefaultInstance.JPCartoon(form)
		fmt.Println(status, err)
		b, _ := json.Marshal(resp)
		log.Info(fmt.Sprintf("JPCartoon resp=%v", string(b)))

		if resp == nil || resp.Code != 10000 {
			return resp.Message, false
		}

		return resp.Data.Image, true
	} else if aiType == 2 {
		resp, status, err := visual.DefaultInstance.ConvertPhoto(form)
		fmt.Println(status, err)
		b, _ := json.Marshal(resp)
		log.Info(fmt.Sprintf("ConvertPhoto resp=%v", string(b)))

		if resp == nil || resp.Code != 10000 {
			return resp.Message, false
		}

		return resp.Data.Image, true
	}

	return "system error", false

}
