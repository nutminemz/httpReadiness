package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"line/health/model"

	"github.com/spf13/viper"
)

// LineResult :: POST to line API
func LineResult(t int, s int, f int, u int64) (e error) {

	client := &http.Client{}
	reqBody := model.HealthRequest{
		TotalWebsites: t,
		Success:       s,
		Failure:       f,
		TotalTime:     u,
	}
	endpoint := viper.GetString("line.report.endpoint")
	token := viper.GetString("line.report.token")
	log.Println("POST ", endpoint)
	log.Printf("BODY %+v\n", reqBody)
	postBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(postBody))
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	log.Println("response code: ", resp.StatusCode)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	log.Printf(newStr)
	return nil
}
