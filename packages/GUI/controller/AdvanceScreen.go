package controller

import (
	"fmt"

	"fyne.io/fyne/v2/widget"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
)

func PlatformSelectionForKeyAddition(s string) {

}

func KeyEntryChange(s string) {

}

func AddKeySubmit() {

}

func ServerUrlChange(url string) {
	model.AddToDb(model.SERVERURL, url)
}

func TranslationUrlChange(url string) {
	model.AddToDb(model.TRANSLATIONURL, url)
}

func ConfigSubmit() {

}

func ChangeProfileKeyAddtion(entry *widget.SelectEntry, options []string) func(s string) {
	return func(s string) {
		var filterOptions = []string{}

		for _, option := range options {
			if isMatch(option, s) {
				filterOptions = append(filterOptions, option)
			}
		}

		entry.SetOptions(filterOptions)
		entry.ActionItem.Show()
		fmt.Println(s)
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
