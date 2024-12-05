package voiceclient

//get voices -> string -> []string
//get voiceId -> string -> string
//convert voice -> line string, filename string -> writes to the file
//new

type IVoiceConversion interface {
	GetVoices(string) []string
	GetVoiceId(string) string
	ConvertVoice(string, string) error
	New() IVoiceConversion
}

var VoiceClientDir = map[string]IVoiceConversion{}

func GetRegistedPlatforms() []string {
	var platforms = []string{}

	for k := range VoiceClientDir {
		platforms = append(platforms, k)
	}

	return platforms
}
