package voiceclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
)

type DanielLabs struct {
}

func init() {
	VoiceClientDir["Daniellabs"] = &DanielLabs{}
	model.AddToDb(model.DANIELLABSURL, "192.168.29.151:7861")
}

func (obj *DanielLabs) New() IVoiceConversion {
	return &DanielLabs{}
}

func (obj *DanielLabs) GetVoices(key string) []string {
	return []string{
		"Daniel-finetuned",
	}
}

func (obj *DanielLabs) GetVoiceId(voice string) string {
	return "Daniel-tuned1"
}

func (obj *DanielLabs) ConvertVoice(line string, filePath string) error {
	hash := addToJobQueueDaniel(line)
	if hash == "" {
		return nil
	}

	err := checkJobDaniel(filePath, hash)
	if err != nil {
		errorhandling.HandleErrorPop(err)
	}

	return nil
}

func addToJobQueueDaniel(line string) string {
	var hash = uuid.New()
	var body = map[string]any{}
	var refFile = map[string]any{
		"path": "https://github.com/jala-R/Voices/raw/refs/heads/main/daniel_voice_sample.mp3",
		"meta": map[string]string{
			"_type": "gradio.FileData",
		},
		"mime_type": "audio/mpeg",
		"orig_name": "daniel_voice_sample.mp3",
		"size":      131019,
		"url":       "https://github.com/jala-R/Voices/raw/refs/heads/main/daniel_voice_sample.mp3",
	}

	var data = []any{
		refFile,
		"It flows through our bodies, radiates in our surroundings, and connects us to the world in ways that often go unseen.",
		line,
		"F5-TTS",
		"true",
		0.1,
		1,
	}

	body["data"] = data
	body["event_data"] = nil
	body["fn_index"] = 0
	body["session_hash"] = hash.String()
	body["trigger_id"] = 6

	bytesData, _ := json.Marshal(body)

	clientUrl := getUrlDaniel()
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/gradio_api/queue/join", clientUrl), bytes.NewReader(bytesData))
	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return ""
	}

	return hash.String()
}

func checkJobDaniel(filename string, hash string) error {
	apiUrl := getUrlDaniel()
	resp, err := http.Get(fmt.Sprintf("http://%s/gradio_api/call/infer/%s", apiUrl, hash))
	if err != nil {
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var urlStr = `http://`

	start := bytes.Index(data, []byte(urlStr))

	if start == -1 {
		return errors.New("audio job failed")
	}

	url := data[start:]

	url = bytes.ReplaceAll(url, []byte(`\\`), []byte(`\`))
	url = bytes.Split(url, []byte(`"`))[0]

	fmt.Println(string(url))
	dataResp, _ := http.Get(string(url))

	file, _ := os.Create(filename)
	file.ReadFrom(dataResp.Body)

	return nil

}

func getUrlDaniel() string {
	val := model.QueryDB(model.DANIELLABSURL)
	url, _ := (val).(string)

	return url
}
