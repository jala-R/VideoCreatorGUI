package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/controller"
)

func advanceTabGUI() fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("Add Key", addKeyTab()),
		container.NewTabItem("Config", configUrls()),
	)
}

func addKeyTab() fyne.CanvasObject {
	platformOptions := widget.NewSelect([]string{"11 labs", "Play HT"}, controller.PlatformSelectionForKeyAddition)
	profileOptions := []string{"abc", "deaf", "sdav", "dbfv"}
	profileSelection := widget.NewSelectEntry(profileOptions)

	profileSelection.OnChanged = controller.ChangeProfileKeyAddtion(profileSelection, profileOptions)

	keyEntry := createSingleEntry("Enter key", controller.KeyEntryChange)

	form := widget.NewForm(
		widget.NewFormItem("Platform", platformOptions),
		widget.NewFormItem("Profile", profileSelection),
		widget.NewFormItem("Key", keyEntry),
	)

	form.OnSubmit = controller.AddKeySubmit

	return form
}

func configUrls() fyne.CanvasObject {

	transUrl := widget.NewEntry()
	transUrl.OnChanged = controller.TranslationUrlChange
	transUrl.SetPlaceHolder("Translation url")

	serverUrl := widget.NewEntry()
	serverUrl.OnChanged = controller.ServerUrlChange
	serverUrl.SetPlaceHolder("Server url")

	form := widget.NewForm(
		widget.NewFormItem("Server url", serverUrl),
		widget.NewFormItem("Transaltion url", transUrl),
	)

	form.OnSubmit = controller.ConfigSubmit

	return form
}