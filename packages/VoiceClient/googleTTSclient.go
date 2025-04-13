package voiceclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
)

//	'voice':{
//	  'languageCode':'en-gb',
//	  'name':'en-GB-Standard-A',
//	  'ssmlGender':'FEMALE'
//	},
type voiceDetails struct {
	LangCode   string `json:"languageCode"`
	Name       string `json:"name"`
	Gender     string `json:"ssmlGender"`
	sampleRate int    `json:"-"`
}

type GoogleLabsClient struct {
	voiceDetail      voiceDetails
	apiKey           string
	selectedVoiceKey []string
	project          string
}

func init() {
	VoiceClientDir["Googlelabs"] = &GoogleLabsClient{}

}

func (obj *GoogleLabsClient) New() IVoiceConversion {
	return &GoogleLabsClient{}
}

func (obj *GoogleLabsClient) GetVoices(key string) []string {
	keySplit := strings.Split(key, " ")
	obj.apiKey = keySplit[1]
	obj.project = keySplit[0]
	url := "https://texttospeech.googleapis.com/v1/voices?language_code=en-US"
	header := http.Header{}
	header["Accept"] = []string{"application/json"}
	header["Content-Type"] = []string{"application/json"}
	header["Authorization"] = []string{"Bearer " + obj.apiKey}
	header["x-goog-user-project"] = []string{obj.project}
	header["Content-Type"] = []string{"application/json", "charset=utf-8"}

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

	var voiceresp = map[string][]VoiceResponse{}

	if err = json.Unmarshal(body, &voiceresp); err != nil {
		errorhandling.HandleErrorPop(fmt.Errorf("Got error while unmarshalling %w", err))
		return nil
	}

	var voiceFound = []string{}

	for _, voice := range voiceresp["voices"] {
		voiceFound = append(voiceFound, fmt.Sprintf("%s %s %s %d", voice.Name, voice.Ssmlgender, voice.LanCode[0], voice.SampleRate))
	}

	return voiceFound
}

type VoiceResponse struct {
	LanCode    []string `json:"languageCodes"`
	Name       string   `json:"name"`
	Ssmlgender string   `json:"ssmlGender"`
	SampleRate int      `json:"naturalSampleRateHertz"`
}

func (obj *GoogleLabsClient) GetVoiceId(voice string) string {
	voiceSplit := strings.Split(voice, " ")

	sampleRate, err := strconv.Atoi(voiceSplit[3])

	if err != nil {
		errorhandling.HandleErrorPop(fmt.Errorf("Error while convertaing samplerate %w", err))
		return ""
	}

	obj.voiceDetail = voiceDetails{
		voiceSplit[2],
		voiceSplit[0],
		voiceSplit[1],
		sampleRate,
	}
	return voice
}

func (obj *GoogleLabsClient) ConvertVoice(line string, filePath string) error {

	url := "https://texttospeech.googleapis.com/v1/text:synthesize"

	header := make(map[string][]string)
	header["Accept"] = []string{"application/json"}
	header["Content-Type"] = []string{"application/json"}
	header["Authorization"] = []string{"Bearer " + obj.apiKey}
	header["x-goog-user-project"] = []string{obj.project}
	header["Content-Type"] = []string{"application/json", "charset=utf-8"}
	body := GetGoogleRequestData(obj.voiceDetail)
	body.Input["text"] = line

	bodyJson, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJson))
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

	return handleGoogleTtsResponse(res.Body, file, obj.voiceDetail.sampleRate)

}

func handleGoogleTtsResponse(body io.Reader, file *os.File, samplerate int) error {
	data, _ := io.ReadAll(body)

	var response = map[string]string{}

	json.Unmarshal(data, &response)

	// fmt.Println(len(response["audioContent"]))

	decodedData, err := base64.StdEncoding.DecodeString(response["audioContent"])
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	return convertMp3ToWav(bytes.NewReader(decodedData), file, samplerate)

}

type TTSRequest struct {
	Input       map[string]string `json:"input"`
	Voice       map[string]string `json:"voice"`
	AudioConfig map[string]string `json:"audioConfig"`
}

func GetGoogleRequestData(voice voiceDetails) TTSRequest {
	var ans = TTSRequest{
		map[string]string{},
		map[string]string{},
		map[string]string{},
	}

	ans.Voice["languageCode"] = voice.LangCode
	ans.Voice["name"] = voice.Name
	ans.Voice["ssmlGender"] = voice.Gender

	ans.AudioConfig["audioEncoding"] = "mp3"

	return ans

}
