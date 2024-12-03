package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

//general
//create projects
//set project name
//set output folder
//set images folder

//audio
//select script
//select languages
//uploaded file -> //english
//german
//spanish
//portugal
//convert audio
//

//advance
//Add voice platform
//Add profile
//Add Keys

func MainScreenGUI(w fyne.Window) fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("General", generalTabGUI(w)),
		container.NewTabItem("Script", scriptTabGUI(w)),
		container.NewTabItem("Advance", advanceTabGUI()),
	)
}
