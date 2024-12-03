package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/view/types"
)

func createSingleEntry(placeholder string, onChange func(string)) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	entry.OnChanged = onChange
	entry.Validator = validation.NewRegexp(`[\S\s]+[\S]+`, "Not Valid ProjectName")
	return entry
}

func selectResourceDialog(resource types.ResourceType) fyne.CanvasObject {

	//path entry
	resourceLablel := widget.NewLabel("")
	resourceLablel.Hide()

	//select button
	browseButton := resource.CreateButton(resourceLablel)

	renameButton := widget.NewButton("Cancel", func() {
		resource.Flush(resourceLablel)
	})

	hBox := container.NewHBox(
		resourceLablel,
		browseButton,
		renameButton,
	)

	return hBox
}
