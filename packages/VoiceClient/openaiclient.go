package voiceclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
)

type OpenAiVoiceClient struct {
	selectedVoiceKey string
	apiKey           string
}

func init() {
	VoiceClientDir["OpenAIVoice"] = &OpenAiVoiceClient{}

}

func (obj *OpenAiVoiceClient) New() IVoiceConversion {
	return &OpenAiVoiceClient{}
}

func (obj *OpenAiVoiceClient) GetVoices(key string) []string {
	obj.apiKey = key

	return []string{
		"echo",
		"fable",
		"alloy",
	}
}

func (obj *OpenAiVoiceClient) GetVoiceId(voice string) string {
	obj.selectedVoiceKey = voice
	return obj.selectedVoiceKey
}

func (obj *OpenAiVoiceClient) ConvertVoice(line string, filePath string) error {

	var body = map[string]string{}

	body["model"] = "gpt-4o-mini-tts"
	body["voice"] = obj.selectedVoiceKey
	body["input"] = line

	body["instructions"] = `Voice: Clear, confident, deep male voice sounding mid-thirties.
	Tone: Highly engaging, dynamic, and narrative. Energetic and enthusiastic, pulling listeners directly into the story like an exciting retelling.
	Emotion: Dynamically mirrors the manhwa's excitement during action, suspense in cliffhangers, intrigue with plot twists, and key character emotions.
	Pacing: Brisk and flowing, faster overall for efficient recapping. Speed varies naturally with content intensity while maintaining strong momentum.
	Clarity: Crisp, articulate pronunciation ensures complex names and plot points are easily understood, even at speed. Uses natural emphasis to highlight key info.
	Pausing: Employs brief, strategic pauses for dramatic impact (reveals, scene breaks, key lines) without disrupting the overall energetic flow.
	`
	body["response_format"] = "mp3"

	data, err := json.Marshal(body)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	// fmt.Println(string(data))

	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/audio/speech",
		bytes.NewReader(data),
	)

	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+obj.apiKey)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		content, _ := io.ReadAll(resp.Body)
		errorhandling.HandleErrorPop(fmt.Errorf("elevenlabs api error: status code : %d with message : %s", resp.StatusCode, content))
		if resp.StatusCode == 429 {
			return errors.New("eleven labs error")
		}
		return nil
	}

	fmt.Println(resp.Status, resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))

	file, err := os.Create(filePath)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return nil
	}
	defer file.Close()

	return convertMp3ToWav(resp.Body, file, 24000)

}
