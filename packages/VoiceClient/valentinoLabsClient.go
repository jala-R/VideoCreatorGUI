package voiceclient

type VelentioLabs struct {
	voiceKey         [][]string
	apiKey           string
	selectedVoiceKey []string
}

func init() {
	VoiceClientDir["valentinolabs"] = &VelentioLabs{}
}

func (obj *VelentioLabs) New() IVoiceConversion {
	return &VelentioLabs{}
}

func (obj *VelentioLabs) GetVoices(key string) []string {
	return nil
}

func (obj *VelentioLabs) GetVoiceId(voice string) string {
	return ""
}

func (obj *VelentioLabs) ConvertVoice(line string, filePath string) error {
	return nil
}
