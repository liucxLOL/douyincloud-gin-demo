package handle_volc

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/volcengine/volc-sdk-golang/service/visual"
)

func GetAIPhoto(imageBase64 string, aiType int) string {
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
		fmt.Println(string(b))

		if resp.Code != 10000 {
			return ""
		}

		return resp.Data.Image
	} else if aiType == 2 {
		resp, status, err := visual.DefaultInstance.ConvertPhoto(form)
		fmt.Println(status, err)
		b, _ := json.Marshal(resp)
		fmt.Println(string(b))

		if resp.Code != 10000 {
			return ""
		}

		return resp.Data.Image
	}

	return ""

}
