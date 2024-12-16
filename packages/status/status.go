package status

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var window fyne.Window

func Register(w fyne.Window) {
	window = w
}

func Pop() {
	info := dialog.NewInformation("Status", "Done", window)
	info.Resize(fyne.NewSize(200, 200))
	info.Show()
}

func ErrorPop(msg string) {
	info := dialog.NewInformation("Status", msg, window)
	info.Resize(fyne.NewSize(200, 200))
	info.Show()
}
