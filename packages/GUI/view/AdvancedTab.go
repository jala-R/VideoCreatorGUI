package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/controller"
	voiceclient "github.com/jala-R/VideoAutomatorGUI/packages/VoiceClient"
)

func advanceTabGUI() fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("Add Key", addKeyTab()),
		container.NewTabItem("Config", configUrls()),
		container.NewTabItem("Key counts", keyCounts()),
	)
}

func keyCounts() fyne.CanvasObject {
	keysDetails := widget.NewMultiLineEntry()
	keysDetails.SetMinRowsVisible(35)
	platformOptions := widget.NewSelect(voiceclient.GetRegistedPlatforms(), controller.OnKeyViewPlatformChange(keysDetails))

	form := widget.NewForm(
		widget.NewFormItem("Platform", platformOptions),
		widget.NewFormItem("", keysDetails),
	)

	return form
}

func addKeyTab() fyne.CanvasObject {
	profileOptions := []string{}
	profileSelection := widget.NewSelectEntry(profileOptions)
	platformOptions := widget.NewSelect(voiceclient.GetRegistedPlatforms(), controller.PlatformSelectionForKeyAddition(profileSelection))

	profileSelection.OnChanged = controller.ChangeProfileKeyAddtion(profileSelection)

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

	valentinoLabsUrl := widget.NewEntry()
	valentinoLabsUrl.SetPlaceHolder("valentino voice url")
	valentinoLabsUrl.OnChanged = controller.ValentinoVoiceUrlChange

	form := widget.NewForm(
		widget.NewFormItem("Server url", serverUrl),
		widget.NewFormItem("Transaltion url", transUrl),
		widget.NewFormItem("Valentino url", valentinoLabsUrl),
	)

	form.OnSubmit = controller.ConfigSubmit

	return form
}
