package voiceclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

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
		return getVoiceOnlyElevenLabs(obj.voiceKey)
	}
	obj.apiKey = key
	url := "https://api.elevenlabs.io/v1/voices"
	header := http.Header{}
	header["Accept"] = []string{"application/json"}
	header["xi-api-key"] = []string{key}
	header["Content-Type"] = []string{"application/json"}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	req.Header = header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	body, _ := io.ReadAll(
		res.Body,
	)

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		errorhandling.HandleErrorPop(fmt.Errorf("voice for 11labs status code %d with message %s and response %s", res.StatusCode, res.Status, string(body)))
		return nil
	}

	obj.voiceKey = getAllVoiceDetailsElevenLabs(body)

	return getVoiceOnlyElevenLabs(obj.voiceKey)
}

func getAllVoiceDetailsElevenLabs(body []byte) [][]string {
	var voiceKey = [][]string{}
	var responseMap = map[string]any{}

	err := json.Unmarshal(body, &responseMap)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	voicesDetail, ok := responseMap["voices"].([]any)
	if !ok {
		errorhandling.HandleErrorPop(fmt.Errorf("getting voice details from eleven labs response json type conversion error"))
		return nil
	}

	for _, t := range voicesDetail {
		temp, _ := t.(map[string]any)
		name, ok := temp["name"].(string)
		if !ok {
			errorhandling.HandleErrorPop(fmt.Errorf("getting voice details from eleven labs response json type conversion name error"))
			return nil
		}
		voiceId, ok := temp["voice_id"].(string)
		if !ok {
			errorhandling.HandleErrorPop(fmt.Errorf("getting voice details from eleven labs response json type conversion name error"))
			return nil
		}

		voiceKey = append(voiceKey, []string{name, voiceId})
	}

	return (voiceKey)

}

func getVoiceOnlyElevenLabs(voiceKey [][]string) []string {
	var ans = []string{}
	for _, temp := range voiceKey {
		ans = append(ans, temp[0])
	}

	slices.Sort(ans)

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
	var selectedVoiceKey = obj.selectedVoiceKey[1]

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s/stream", selectedVoiceKey)

	headers := make(map[string][]string)
	headers["Accept"] = []string{"application/json"}
	headers["xi-api-key"] = []string{obj.apiKey}

	body := GetRequestData()
	body.Text = line

	bodyJson, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJson))
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	req.Header = headers

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		content, _ := io.ReadAll(res.Body)
		errorhandling.HandleErrorPop(fmt.Errorf("elevenlabs api error: status code : %d with message : %s", res.StatusCode, content))
		if res.StatusCode == 401 {
			return errors.New("eleven labs error")
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	defer file.Close()

	return convertMp3ToWav(res.Body, file, 44100)

}

type RequestData struct {
	Text          string        `json:"text"`
	Model_id      string        `json:"model_id"`
	VoiceSettings VoiceSettings `json:"voice_settings"`
}

type VoiceSettings struct {
	Stability         float64 `json:"stability"`
	Similarity_boost  float64 `json:"similarity_boost"`
	Style             float64 `json:"style"`
	Use_speaker_boost bool    `json:"use_speaker_boost"`
}

func GetRequestData() RequestData {
	data := RequestData{}
	data.Model_id = "eleven_multilingual_v2"
	data.VoiceSettings.Stability = 0.5
	data.VoiceSettings.Similarity_boost = 0.8
	data.VoiceSettings.Style = 0.0
	data.VoiceSettings.Use_speaker_boost = true
	return data
}
