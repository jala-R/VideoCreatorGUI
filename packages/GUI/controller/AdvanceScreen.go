package controller

import (
	"bufio"
	"fmt"
	"strings"

	"fyne.io/fyne/v2/widget"
	apiclient "github.com/jala-R/VideoAutomatorGUI/packages/ApiClient"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
	"github.com/jala-R/VideoAutomatorGUI/packages/status"
)

func PlatformSelectionForKeyAddition(profile *widget.SelectEntry) func(s string) {
	return func(s string) {
		temp := apiclient.ListProfileOnPlatform(s)
		profile.SetOptions(temp)
		model.AddToDb(model.PROFILES, temp)
		model.AddToDb(model.PLATFORMADD, s)
	}
}

func KeyEntryChange(s string) {
	model.AddToDb(model.KEYADD, s)
}

func AddKeySubmit() {
	var keys = []string{model.PLATFORMADD, model.PROFILEADD, model.KEYADD}
	var values = [3]string{}

	for i, key := range keys {
		val := model.QueryDB(key)
		if val == nil {
			return
		}

		temp := val.(string)
		if temp == "" {
			return
		}
		values[i] = temp
	}

	for _, key := range keys {
		model.AddToDb(key, nil)
	}

	apiclient.AddKey(values[0], values[1], values[2])

	status.Pop()

}

func ServerUrlChange(url string) {
	model.AddToDb(model.APIURL, url)
}

func TranslationUrlChange(url string) {
	model.AddToDb(model.TRANSLATIONURL, url)
}

func ValentinoVoiceUrlChange(url string) {
	model.AddToDb(model.VALENTINOVOICEURL, url)
}

func DanielVoiceUrlChange(url string) {
	model.AddToDb(model.DANIELLABSURL, url)
}

func KokoroVoiceUrlChange(url string) {
	model.AddToDb(model.KOKOROLABSURL, url)
}

func ConfigSubmit() {

}

func ChangeProfileKeyAddtion(entry *widget.SelectEntry) func(s string) {
	return func(s string) {
		val := model.QueryDB(model.PROFILES)
		options := []string{}
		if val != nil {
			options, _ = val.([]string)
		}

		var filterOptions = []string{}

		for _, option := range options {
			if isMatch(option, s) {
				filterOptions = append(filterOptions, option)
			}
		}

		entry.SetOptions(filterOptions)
		entry.ActionItem.Show()
		model.AddToDb(model.PROFILEADD, s)
	}
}

func isMatch(str string, pattern string) bool {
	temp := pattern + "^" + str
	lps := make([]int, len(temp))

	left := 0
	right := 1

	for right < len(lps) {
		if temp[left] == temp[right] {
			left++
			lps[right] = left
			right++
		} else {
			if left == 0 {
				right++
			} else {
				left = lps[left-1]
			}
		}
	}

	for _, match := range lps {
		if match == len(pattern) {
			return true
		}
	}

	return false
}

func OnKeyViewPlatformChange(entry *widget.Entry) func(string) {
	return func(platform string) {

		//presist the platform
		model.AddToDb(model.PLATFORMFORKEYSVIEW, platform)

		//make api call to server and get the profile to key count
		keyCounts := apiclient.GetKeyDetailsForPlatform(platform)
		var textToShow = make([]byte, 0, 100)

		for k, v := range keyCounts {
			textToShow = append(textToShow, fmt.Sprintf("%s : %d\n", k, v)...)
		}

		entry.SetText(string(textToShow))
	}
}

func SaveProfiles(entry *widget.Entry) func() {
	return func() {
		text := entry.Text

		ip := bufio.NewReader(strings.NewReader(text))
		profiles := []string{}
		for {
			line, _ := ip.ReadString('\n')
			if len(line) == 0 {
				break
			}
			line = line[:len(line)-1]

			prof := strings.Split(line, ":")[0]

			if prof[len(prof)-1] == ' ' {
				prof = prof[:len(prof)-1]
			}

			profiles = append(profiles, prof)
		}

		apiclient.PresistProfiles(profiles)

	}
}
