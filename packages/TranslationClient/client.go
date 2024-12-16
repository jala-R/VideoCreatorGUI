package translationclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
)

func init() {
	model.AddToDb(model.TRANSLATIONAPIURL, "192.168.29.151:8080")
}

//get the model
//add prefix

var langVsModel = map[string]string{
	model.LOCALES[3]: "Helsinki-NLP/opus-mt-tc-big-en-pt",
	model.LOCALES[1]: "Helsinki-NLP/opus-mt-en-de",
	model.LOCALES[2]: "Helsinki-NLP/opus-mt-en-es",
}

func TranslateSentence(line, locale string) string {
	if locale == model.LOCALES[3] {
		line = ">>por<< " + line
	}

	var body = map[string]string{}
	body["content"] = line
	body["model"] = langVsModel[locale]

	requestBody, err := json.Marshal(body)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	url := getUrl()
	if url == "" {
		errorhandling.HandleErrorPop(errors.New("no url set in db"))
		return ""
	}

	fullUrl := fmt.Sprintf("http://%s/translation", url)

	req, err := http.NewRequest(
		"POST",
		fullUrl,
		bytes.NewReader(requestBody),
	)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	var respMap = map[string]string{}
	err = json.Unmarshal(data, &respMap)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	return respMap["received_data"]
}
