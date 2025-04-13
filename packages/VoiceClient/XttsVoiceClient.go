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

type XttsLabs struct {
	selectedVoice string
}

func init() {
	VoiceClientDir["XttsLabs"] = &XttsLabs{}
	model.AddToDb(model.XTTSLABSURL, "127.0.0.1:8000")
}

func (obj *XttsLabs) New() IVoiceConversion {
	return &XttsLabs{}
}

func (obj *XttsLabs) GetVoices(key string) []string {
	return []string{
		"Claribel Dervla", "Daisy Studious", "Gracie Wise", "Tammie Ema", "Alison Dietlinde", "Ana Florence", "Annmarie Nele", "Asya Anara", "Brenda Stern", "Gitta Nikolina", "Henriette Usha", "Sofia Hellen", "Tammy Grit", "Tanja Adelina", "Vjollca Johnnie", "Andrew Chipper", "Badr Odhiambo", "Dionisio Schuyler", "Royston Min", "Viktor Eka", "Abrahan Mack", "Adde Michal", "Baldur Sanjin", "Craig Gutsy", "Damien Black", "Gilberto Mathias", "Ilkin Urbano", "Kazuhiko Atallah", "Ludvig Milivoj", "Suad Qasim", "Torcull Diarmuid", "Viktor Menelaos", "Zacharie Aimilios", "Nova Hogarth", "Maja Ruoho", "Uta Obando", "Lidiya Szekeres", "Chandra MacFarland", "Szofi Granger", "Camilla Holmström", "Lilya Stainthorpe", "Zofija Kendrick", "Narelle Moon", "Barbora MacLean", "Alexandra Hisakawa", "Alma María", "Rosemary Okafor", "Ige Behringer", "Filip Traverse", "Damjan Chapman", "Wulf Carlevaro", "Aaron Dreschner", "Kumar Dahl", "Eugenio Mataracı", "Ferran Simen", "Xavier Hayasaka", "Luis Moray", "Marcos Rudaski",
	}
}

func (obj *XttsLabs) GetVoiceId(voice string) string {
	obj.selectedVoice = voice
	return voice
}

func (obj *XttsLabs) ConvertVoice(line string, filePath string) error {
	url := obj.getUrl()
	fullUrl := fmt.Sprintf("http://%s/CreateAudio", url)

	reqBody := map[string]string{
		"line":            line,
		"output_location": filePath,
		"voice":           obj.selectedVoice,
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

func (obj *XttsLabs) getUrl() string {
	val := model.QueryDB(model.XTTSLABSURL)
	url, _ := (val).(string)

	return url
}
