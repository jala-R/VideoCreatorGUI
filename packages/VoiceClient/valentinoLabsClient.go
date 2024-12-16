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

type VelentioLabs struct {
}

func init() {
	VoiceClientDir["valentinolabs"] = &VelentioLabs{}
	model.AddToDb(model.VALENTINOVOICEURL, "192.168.29.151:7860")
}

func (obj *VelentioLabs) New() IVoiceConversion {
	return &VelentioLabs{}
}

func (obj *VelentioLabs) GetVoices(key string) []string {
	return []string{
		"Valentino-finetuned",
	}
}

func (obj *VelentioLabs) GetVoiceId(voice string) string {
	return "valentino-tuned1"
}

func (obj *VelentioLabs) ConvertVoice(line string, filePath string) error {
	hash := addToJobQueue(line)
	if hash == "" {
		return nil
	}

	err := checkJob(filePath, hash)
	if err != nil {
		errorhandling.HandleErrorPop(err)
	}

	return nil
}

func addToJobQueue(line string) string {
	var hash = uuid.New()
	var body = map[string]any{}
	var refFile = map[string]any{
		"path": "https://github.com/jala-R/Voices/raw/refs/heads/main/valentino_sample.wav",
		"meta": map[string]string{
			"_type": "gradio.FileData",
		},
		"mime_type": "audio/wav",
		"orig_name": "valentino_sample.wav",
		"size":      921644,
		"url":       "https://github.com/jala-R/Voices/raw/refs/heads/main/valentino_sample.wav",
	}

	var data = []any{
		refFile,
		"Now, a dramatic change is underway. Those who sought your downfall are now seeing their own schemes backfire spectacularly, as if thrown back by a boomerang.",
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
	body["trigger_id"] = 7

	bytesData, _ := json.Marshal(body)

	clientUrl := getUrl()
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

func checkJob(filename string, hash string) error {
	apiUrl := getUrl()
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

func getUrl() string {
	val := model.QueryDB(model.VALENTINOVOICEURL)
	url, _ := (val).(string)

	return url
}
