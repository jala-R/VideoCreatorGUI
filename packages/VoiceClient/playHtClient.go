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
	"strings"

	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
)

type PlayHTClient struct {
	voiceKey         [][]string //name id
	apiKey           []string
	selectedVoiceKey []string
}

func init() {
	VoiceClientDir["playHt"] = &PlayHTClient{}
}

func (obj *PlayHTClient) New() IVoiceConversion {
	return &PlayHTClient{}
}

func parseKey(key string) []string {
	userIdKey := strings.Split(key, "-")
	return userIdKey
}

func (obj *PlayHTClient) GetVoices(key string) []string {
	if key == "" {
		errorhandling.HandleErrorPop(errors.New("no key in selected profile"))
		return nil
	}
	if obj.voiceKey != nil {
		return getVoiceOnlyPlayHT(obj.voiceKey)
	}
	obj.apiKey = parseKey(key)
	url := "https://api.play.ht/api/v2/voices"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	req.Header["X-USER-ID"] = []string{obj.apiKey[0]}
	req.Header["AUTHORIZATION"] = []string{obj.apiKey[1]}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	obj.voiceKey = getAllVoiceDetailsPlayHT(body)
	return getVoiceOnlyPlayHT(obj.voiceKey)

}

func getVoiceOnlyPlayHT(voiceKey [][]string) []string {
	var ans = []string{}
	for _, temp := range voiceKey {
		ans = append(ans, temp[0])
	}

	slices.Sort(ans)

	return ans
}

func getAllVoiceDetailsPlayHT(body []byte) [][]string {
	var voiceKey = [][]string{}
	var temp []any

	err := json.Unmarshal(body, &temp)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	for _, voice := range temp {
		voiceDetails, ok := voice.(map[string]any)
		if !ok {
			errorhandling.HandleErrorPop(fmt.Errorf("error at play ht array type casting"))
			return nil
		}

		val := voiceDetails["name"]
		if val == nil {
			continue
		}
		name, ok := val.(string)
		if !ok {
			continue
		}

		val = voiceDetails["voice_engine"]
		if val == nil {
			continue
		}
		voiceEngine, ok := val.(string)
		if !ok {
			continue
		}

		val = voiceDetails["language"]
		if val == nil {
			continue
		}
		language, ok := val.(string)
		if !ok {
			continue
		}

		fullname := fmt.Sprintf("%s-%s-%s", name, voiceEngine, language)

		val = voiceDetails["id"]
		if val == nil {
			continue
		}
		id, ok := val.(string)
		if !ok {
			continue
		}

		voiceKey = append(voiceKey, []string{fullname, id})

	}

	//manual voice upload for mutilungulag

	//spanish
	voiceKey = append(voiceKey, []string{
		"xavi-spanish",
		"s3://voice-cloning-zero-shot/36328a44-5c42-4a35-a9a1-b45596a56c88/original/manifest.json",
	})

	//spanish
	voiceKey = append(voiceKey, []string{
		"Patricia Narrative",
		"s3://voice-cloning-zero-shot/5694d5e5-2dfe-4440-8cc8-e2a69c3e7560/original/manifest.json",
	})

	//german2
	voiceKey = append(voiceKey, []string{
		"illas-german",
		"s3://voice-cloning-zero-shot/f78a1dc3-6533-4967-a0d2-88e13894a45a/original/manifest.json",
	})

	//portugal
	voiceKey = append(voiceKey, []string{
		"Jacile Narrative",
		"s3://voice-cloning-zero-shot/6d093315-8da6-4fd8-a4b9-5446b43ff4c7/original/manifest.json",
	})

	//portugal
	voiceKey = append(voiceKey, []string{
		"jorge-portugal",
		"s3://voice-cloning-zero-shot/ec8095bd-bbab-4229-8527-0b0ead293823/original/manifest.json",
	})

	return voiceKey

}

func (obj *PlayHTClient) GetVoiceId(voice string) string {
	for _, temp := range obj.voiceKey {
		if voice == temp[0] {
			obj.selectedVoiceKey = make([]string, 2)
			obj.selectedVoiceKey[0] = temp[0]
			obj.selectedVoiceKey[1] = temp[1]
			return obj.selectedVoiceKey[0] + "-" + obj.selectedVoiceKey[1]
		}
	}
	return ""
}

func (obj *PlayHTClient) ConvertVoice(line string, filePath string) error {
	voiceId := obj.selectedVoiceKey[1]

	var body = map[string]any{}
	body["text"] = line
	body["voice_engine"] = "PlayDialog"
	body["voice"] = voiceId
	body["output_format"] = "mp3"
	body["sample_rate"] = 44100

	bodyData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.play.ht/api/v2/tts/stream", bytes.NewReader(bodyData))
	req.Header.Add("X-USER-ID", obj.apiKey[0])
	req.Header.Add("AUTHORIZATION", obj.apiKey[1])
	req.Header.Add("accept", "audio/mpeg")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode == 403 {
		return errors.New("credits over")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		data, _ := io.ReadAll(res.Body)
		errorhandling.HandleError(fmt.Errorf("got error from playht client with code %d, message:%s", res.StatusCode, string(data)))
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		errorhandling.HandleError(err)
		return nil
	}
	defer file.Close()

	convertMp3ToWav(res.Body, file)
	// file.ReadFrom(res.Body)

	return err
}
