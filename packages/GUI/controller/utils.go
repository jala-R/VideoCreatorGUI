package controller

import (
	"reflect"

	"fyne.io/fyne/v2/widget"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
)

func IsEmpty(val any) bool {
	if val == nil {
		return true
	}
	reflectVal := reflect.ValueOf(val)
	return reflectVal.IsZero()
}

func RegisterEntry(key string, widget *widget.Entry) {
	model.AddToDb(key, widget)
}

func GetLocaleOutputEntry(key string) *widget.Entry {
	translatedKey := key + model.LOCALESCRIPTSUFFIX
	val := model.QueryDB(translatedKey)
	if val == nil {
		return nil
	}
	entry, ok := val.(*widget.Entry)
	if !ok {
		return nil
	}

	return entry

}
