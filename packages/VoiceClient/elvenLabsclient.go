package voiceclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
)

type ElevnLabsClient struct {
	voiceKey         [][]string
	apiKey           string
	selectedVoiceKey []string
}

func init() {
	VoiceClientDir["11labs"] = &ElevnLabsClient{}

}

func (obj *ElevnLabsClient) New() IVoiceConversion {
	return &ElevnLabsClient{}
}

func (obj *ElevnLabsClient) GetVoices(key string) []string {
	if obj.voiceKey != nil {
		return getVoiceOnly(obj.voiceKey)
	}
	obj.apiKey = key
	url := "https://api.elevenlabs.io/v1/voices"
	header := http.Header{}
	header["Accept"] = []string{"application/json"}
	header["xi-api-key"] = []string{key}
	header["Content-Type"] = []string{"application/json"}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	req.Header = header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}

	body, _ := io.ReadAll(
		res.Body,
	)

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		errorhandling.HandleError(fmt.Errorf("voice for 11labs status code %d with message %s and response %s", res.StatusCode, res.Status, string(body)))
		return nil
	}

	obj.voiceKey = getAllVoiceDetails(body)

	return getVoiceOnly(obj.voiceKey)
}

func getAllVoiceDetails(body []byte) [][]string {
	var voiceKey = [][]string{}
	var responseMap = map[string]any{}

	err := json.Unmarshal(body, &responseMap)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}
	voicesDetail, ok := responseMap["voices"].([]any)
	if !ok {
		errorhandling.HandleError(fmt.Errorf("getting voice details from eleven labs response json type conversion error"))
		return nil
	}

	for _, t := range voicesDetail {
		temp, _ := t.(map[string]any)
		name, ok := temp["name"].(string)
		if !ok {
			errorhandling.HandleError(fmt.Errorf("getting voice details from eleven labs response json type conversion name error"))
			return nil
		}
		voiceId, ok := temp["voice_id"].(string)
		if !ok {
			errorhandling.HandleError(fmt.Errorf("getting voice details from eleven labs response json type conversion name error"))
			return nil
		}

		voiceKey = append(voiceKey, []string{name, voiceId})
	}

	return (voiceKey)

}

func getVoiceOnly(voiceKey [][]string) []string {
	var ans = []string{}
	for _, temp := range voiceKey {
		ans = append(ans, temp[0])
	}

	return ans
}

func (obj *ElevnLabsClient) GetVoiceId(voice string) string {
	for _, temp := range obj.voiceKey {
		if voice == temp[0] {
			obj.selectedVoiceKey = make([]string, 2)
			obj.selectedVoiceKey[0] = temp[0]
			obj.selectedVoiceKey[1] = temp[1]
			return obj.selectedVoiceKey[1]
		}
	}
	return ""
}

func (obj *ElevnLabsClient) ConvertVoice(line string, filePath string) error {
	fmt.Println(obj.selectedVoiceKey)
	return nil
}
