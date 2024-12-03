package model

import (
	"fmt"

	"fyne.io/fyne/v2/widget"
)

type Database map[string]any

var db Database

func init() {
	db = Database{}
}

func AddToDb(key string, value any) {
	db[key] = value
}

func QueryDB(key string) any {
	return db[key]
}

func PrintDb() {
	fmt.Println(db)
}

func GetEntry(key string) *widget.Entry {
	val := QueryDB(key)
	if val == nil {
		return nil
	}

	entry, ok := val.(*widget.Entry)

	if !ok {
		return nil
	}

	return entry
}
