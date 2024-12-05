package translationclient

import "github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"

func getUrl() string {
	val := model.QueryDB(model.TRANSLATIONAPIURL)
	if val == nil {
		return ""
	}

	url, ok := val.(string)
	if !ok {
		return ""
	}

	return url
}
