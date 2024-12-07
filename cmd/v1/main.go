package main

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/view"
	"github.com/jala-R/VideoAutomatorGUI/packages/status"
)

func main() {
	defer panicLog()
	a := app.New()
	w := a.NewWindow("Video Automator")
	status.Register(w)

	w.SetContent(view.MainScreenGUI(w))

	w.Resize(fyne.NewSize(1080, 1920))
	w.ShowAndRun()
}

func panicLog() {
	panicMessage := recover()
	err, ok := panicMessage.(error)
	if !ok {
		err = errors.New("panic with non error type")
	}
	errorhandling.HandleError(err)
}
