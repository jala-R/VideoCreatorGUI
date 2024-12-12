package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	apiclient "github.com/jala-R/VideoAutomatorGUI/packages/ApiClient"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/controller"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/view/types"
	voiceclient "github.com/jala-R/VideoAutomatorGUI/packages/VoiceClient"
)

//Input
//select script done
//script page
//select list of language
//select outputFolder done

//output
//eng  //pt //sp

func scriptTabGUI(w fyne.Window) fyne.CanvasObject {

	scriptTabs := container.NewAppTabs(
		container.NewTabItem("Input", inputScriptGUI(w)),
		container.NewTabItem("Processed", ProcessedScript(w)),
	)
	return scriptTabs
}

func ProcessedScript(w fyne.Window) fyne.CanvasObject {

	items := []*container.TabItem{}

	for _, locale := range model.LOCALES {
		items = append(items, container.NewTabItem(locale, createProcessedScriptPage(locale)))
	}
	container := container.NewAppTabs(
		items...,
	)

	return container
}

func createProcessedScriptPage(locale string) fyne.CanvasObject {
	var (
		platform      string
		profileOption string
		voice         string
	)

	script := widget.NewMultiLineEntry()
	script.SetMinRowsVisible(35)
	var profile = []string{}

	controller.RegisterEntryVsLocale(locale, script)

	voices := widget.NewSelect([]string{}, func(s string) {
		voice = s
	})

	voiceProfile := widget.NewSelect(profile, func(s string) {
		profileOption = s
		key := apiclient.GetKey(platform, profileOption)
		voiceClientInst := voiceclient.VoiceClientDir[platform].New()
		voices.Options = voiceClientInst.GetVoices(key)

	})
	statusLabel := widget.NewLabel("Not started")

	form := widget.NewForm(
		widget.NewFormItem("", script),
		widget.NewFormItem("Voice Platform", widget.NewSelect(voiceclient.GetRegistedPlatforms(), func(s string) {
			platform = s
			//apiclinet call to get profiles in his platform
			voiceProfile.Options = apiclient.ListProfileOnPlatform(platform)

		})),
		widget.NewFormItem("Voice Profile", voiceProfile),
		widget.NewFormItem("Voices", voices),
		widget.NewFormItem("Status", statusLabel),
	)

	form.OnSubmit = controller.ConvertVoice(script, locale, &platform, &profileOption, &voice, statusLabel)
	return form
}

func inputScriptGUI(w fyne.Window) fyne.CanvasObject {

	scriptContent := widget.NewMultiLineEntry()
	scriptContent.SetMinRowsVisible(35)

	form := widget.NewForm(
		widget.NewFormItem("Select Script", selectResourceDialog(types.NewFileType("Select file", w, controller.ScriptFileHandler, ".txt"))),
		widget.NewFormItem("", scriptContent),
		widget.NewFormItem("Dont split para", widget.NewCheck("True", controller.SetStrict16WordsPerPara)),
		widget.NewFormItem("Select Audio Output folder", selectResourceDialog(types.NewFolderType("Select", w, controller.AudioLocationHandler))),
		widget.NewFormItem("Locales", widget.NewCheckGroup(
			model.LOCALES[1:],
			controller.SelectLocaleHandler,
		)),
	)

	controller.RegisterEntry(model.INPUTSCRIPTWIDGET, scriptContent)

	form.OnSubmit = controller.ScriptInputSubmit(w)

	return form
}
