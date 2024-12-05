package voiceclient

type PlayHTClient struct {
	voiceKey         [][]string
	apiKey           string
	selectedVoiceKey []string
}

func init() {
	VoiceClientDir["playHt"] = &PlayHTClient{}
}

func (obj *PlayHTClient) New() IVoiceConversion {
	return &PlayHTClient{}
}

func (obj *PlayHTClient) GetVoices(key string) []string {
	return nil
}

func (obj *PlayHTClient) GetVoiceId(voice string) string {
	return ""
}

func (obj *PlayHTClient) ConvertVoice(line string, filePath string) error {
	return nil
}
