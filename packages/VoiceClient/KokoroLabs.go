package voiceclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
)

type KokoroLabs struct {
	selectedVoice string
}

func init() {
	VoiceClientDir["KokoroLabs"] = &KokoroLabs{}
	model.AddToDb(model.KOKOROLABSURL, "127.0.0.1:4000")
}

func (obj *KokoroLabs) New() IVoiceConversion {
	return &KokoroLabs{}
}

func (obj *KokoroLabs) GetVoices(key string) []string {
	return []string{
		"af_alloy",
		"af_aoede",
		"af_bella",
		"af_heart",
		"af_jessica",
		"af_kore",
		"af_nicole",
		"af_nova",
		"af_river",
		"af_sarah",
		"af_sky",
		"am_adam",
		"am_echo",
		"am_eric",
		"am_fenrir",
		"am_liam",
		"am_michael",
		"am_onyx",
		"am_puck",
		"am_santa",
		"bf_alice",
		"bf_emma",
		"bf_isabella",
		"bf_lily",
		"bm_daniel",
		"bm_fable",
		"bm_george",
		"bm_lewis",
		"ef_dora",
		"em_alex",
		"em_santa",
		"ff_siwis",
		"hf_alpha",
		"hf_beta",
		"hm_omega",
		"hm_psi",
		"if_sara",
		"im_nicola",
		"jf_alpha",
		"jf_gongitsune",
		"jf_nezumi",
		"jf_tebukuro",
		"jm_kumo",
		"pf_dora",
		"pm_alex",
		"pm_santa",
		"zf_xiaobei",
		"zf_xiaoni",
		"zf_xiaoxiao",
		"zf_xiaoyi",
	}
}

func (obj *KokoroLabs) GetVoiceId(voice string) string {
	obj.selectedVoice = voice
	return voice
}

func (obj *KokoroLabs) ConvertVoice(line string, filePath string) error {
	url := obj.getUrl()
	fullUrl := fmt.Sprintf("http://%s/convertVoice", url)

	reqBody := map[string]string{
		"content":     line,
		"output_file": filePath,
		"voice":       obj.selectedVoice,
	}

	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(jsonReqBody))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(string(msg))
	}

	return nil

}

// func addToJobQueue(line string) string {
// 	var hash = uuid.New()
// 	var body = map[string]any{}
// 	var refFile = map[string]any{
// 		"path": "https://github.com/jala-R/Voices/raw/refs/heads/main/valentino_sample.wav",
// 		"meta": map[string]string{
// 			"_type": "gradio.FileData",
// 		},
// 		"mime_type": "audio/wav",
// 		"orig_name": "valentino_sample.wav",
// 		"size":      921644,
// 		"url":       "https://github.com/jala-R/Voices/raw/refs/heads/main/valentino_sample.wav",
// 	}

// 	var data = []any{
// 		refFile,
// 		"Now, a dramatic change is underway. Those who sought your downfall are now seeing their own schemes backfire spectacularly, as if thrown back by a boomerang.",
// 		line,
// 		"F5-TTS",
// 		"true",
// 		0.1,
// 		1,
// 	}

// 	body["data"] = data
// 	body["event_data"] = nil
// 	body["fn_index"] = 0
// 	body["session_hash"] = hash.String()
// 	body["trigger_id"] = 7

// 	bytesData, _ := json.Marshal(body)

// 	clientUrl := getUrl()
// 	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/gradio_api/queue/join", clientUrl), bytes.NewReader(bytesData))
// 	req.Header["Content-Type"] = []string{"application/json"}

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		errorhandling.HandleErrorPop(err)
// 		return ""
// 	}

// 	_, err = io.ReadAll(resp.Body)
// 	if err != nil {
// 		errorhandling.HandleErrorPop(err)
// 		return ""
// 	}

// 	return hash.String()
// }

// func checkJob(filename string, hash string) error {
// 	apiUrl := getUrl()
// 	resp, err := http.Get(fmt.Sprintf("http://%s/gradio_api/call/infer/%s", apiUrl, hash))
// 	if err != nil {
// 		return err
// 	}

// 	data, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	var urlStr = `http://`

// 	start := bytes.Index(data, []byte(urlStr))

// 	if start == -1 {
// 		return errors.New("audio job failed")
// 	}

// 	url := data[start:]

// 	url = bytes.ReplaceAll(url, []byte(`\\`), []byte(`\`))
// 	url = bytes.Split(url, []byte(`"`))[0]

// 	fmt.Println(string(url))
// 	dataResp, _ := http.Get(string(url))

// 	file, _ := os.Create(filename)
// 	file.ReadFrom(dataResp.Body)

// 	return nil

// }

func (obj *KokoroLabs) getUrl() string {
	val := model.QueryDB(model.KOKOROLABSURL)
	url, _ := (val).(string)

	return url
}
