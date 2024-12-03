package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/controller"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/view/types"
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
		voices.Options = []string{
			s + "voice 1",
			s + "voice 2",
			s + "voice 3",
		}
		profileOption = s
	})

	form := widget.NewForm(
		widget.NewFormItem("", script),
		widget.NewFormItem("Voice Platform", widget.NewSelect([]string{"11 labs", "play HT"}, func(s string) {
			voiceProfile.Options = []string{
				s + " profile 1",
				s + " profile 2",
				s + " profile 3",
			}
			platform = s
		})),
		widget.NewFormItem("Voice Profile", voiceProfile),
		widget.NewFormItem("Voices", voices),
	)

	form.OnSubmit = controller.ConvertVoice(script, locale, platform, profileOption, voice)
	return form
}

func inputScriptGUI(w fyne.Window) fyne.CanvasObject {

	scriptContent := widget.NewMultiLineEntry()
	scriptContent.SetMinRowsVisible(35)

	form := widget.NewForm(
		widget.NewFormItem("Select Script", selectResourceDialog(types.NewFileType("Select file", w, controller.ScriptFileHandler, ".txt"))),
		widget.NewFormItem("", scriptContent),
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
