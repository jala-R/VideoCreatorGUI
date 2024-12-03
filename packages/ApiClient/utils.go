package apiclient

import (
	"errors"
	"os"
	"strings"

	"github.com/go-audio/wav"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
)

var imageExt = []string{
	"png",
	"jpg",
	"jpeg",
	"webp",
}

func isImageExt(ext string) bool {
	for _, ex := range imageExt {
		if ex == ext {
			return true
		}
	}

	return false
}

func getExtension(filename string) string {
	splitted := strings.Split(filename, ".")
	return splitted[len(splitted)-1]
}

func getWavDuration(filename string) (float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	decoder := wav.NewDecoder(file)

	if !decoder.IsValidFile() {
		return 0, errors.New("Invalid wav file :" + filename)
	}

	dur, err := decoder.Duration()
	if err != nil {
		return 0, err
	}

	return dur.Seconds(), nil
}

func getUrl() string {
	val := model.QueryDB(model.APIURL)
	if val == nil {
		return ""
	}

	url, ok := val.(string)
	if !ok {
		return ""
	}

	return url
}
