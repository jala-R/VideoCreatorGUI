package controller

import (
	"errors"
	"fmt"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	serverclient "github.com/jala-R/VideoAutomatorGUI/packages/ApiClient"
	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
	translationclient "github.com/jala-R/VideoAutomatorGUI/packages/TranslationClient"
	utils "github.com/jala-R/VideoAutomatorGUI/packages/Utils"
)

func ScriptFileHandler(fileLocation string) {
	model.AddToDb(model.SCRIPTFILE, fileLocation)
	//get entry from db
	scriptEntry := model.GetEntry(model.INPUTSCRIPTWIDGET)
	var scriptContent string

	//populate the readed file
	if fileLocation != "" {
		file, err := os.Open(fileLocation)
		if err != nil {
			errorhandling.HandleError(err)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			errorhandling.HandleError(err)
			return
		}

		scriptContent = string(content)
	}
	scriptEntry.SetText(scriptContent)
}

func AudioLocationHandler(outFolder string) {
	model.AddToDb(model.AUDIOOUTPUTFOLDER, outFolder)
}

func SelectLocaleHandler(loclaes []string) {
	model.AddToDb(model.LOCALEOPTIONS, loclaes)
}

func ValidateScriptScreenFeilds() error {
	if IsEmpty(model.QueryDB(model.SCRIPTFILE)) {
		return errors.New("select script file")
	}

	if IsEmpty(model.QueryDB(model.AUDIOOUTPUTFOLDER)) {
		return errors.New("select audio output path")
	}

	scriptEntry := model.GetEntry(model.INPUTSCRIPTWIDGET)

	if IsEmpty(scriptEntry.Text) {
		return errors.New("no script found")
	}

	return nil
}

func ScriptInputSubmit(w fyne.Window) func() {
	return func() {
		err := ValidateScriptScreenFeilds()

		if err != nil {
			errorhandling.HandleError(err)
			return
		}

		//start an go routine for transalation
		go TriggerTranslation()

		info := dialog.NewInformation("Status", "Started translation", w)
		info.Resize(fyne.NewSize(200, 200))
		info.Show()

	}
}

func ConvertVoice(script *widget.Entry, locale string, platform string, profile string, voice string) func() {
	return func() {

	}
}

func RegisterEntryVsLocale(locale string, script *widget.Entry) {
	localeKey := locale + model.LOCALESCRIPTSUFFIX
	model.AddToDb(localeKey, script)
}

func TriggerTranslation() {
	var selectedLocaleVsEntry = map[string]*widget.Entry{}

	//get locale to do + eng
	var val = model.QueryDB(model.LOCALEOPTIONS)
	if val != nil {
		selectedLocales, ok := val.([]string)
		if !ok {
			errorhandling.HandleError(errors.New("type casting failed for selected locale"))
			return
		}

		fmt.Println(selectedLocales)
		for _, locale := range selectedLocales {
			selectedLocaleVsEntry[locale] = GetLocaleOutputEntry(locale)
		}
	}

	var englishScriptEntry = GetLocaleOutputEntry(model.LOCALES[0])
	selectedLocaleVsEntry[model.LOCALES[0]] = englishScriptEntry

	//set all locale as not selected
	for _, locale := range model.LOCALES {
		GetLocaleOutputEntry(locale).SetText("Not Selected")
	}

	fmt.Println(selectedLocaleVsEntry)
	//set the status as started translating
	for _, entry := range selectedLocaleVsEntry {
		entry.SetText("Started Translating")
	}

	//make backendcall to parse the english
	inputEntry := model.GetEntry(model.INPUTSCRIPTWIDGET)
	if inputEntry == nil {
		errorhandling.HandleError(errors.New("no script given"))
		return
	}
	optimizedScript := serverclient.OptimizeScript(inputEntry.Text)
	englishScriptEntry.SetText(optimizedScript)
	delete(selectedLocaleVsEntry, model.LOCALES[0])

	parsedScript := utils.MarshallScript(optimizedScript)

	for locale := range selectedLocaleVsEntry {
		go TranslateEnglishToLocale(locale, parsedScript)
	}

	//take the parsed and translate each sentences set progess in entry
	//set the trnaslated script once tanslation is completed

}

func TranslateEnglishToLocale(locale string, script [][]string) {
	var totalSenetences = 0
	var translatedScript = [][]string{}

	fmt.Println(script)

	for _, para := range script {
		totalSenetences += len(para)
	}

	var currentSenetence = 1
	var entry = GetLocaleOutputEntry(locale)
	//start translating sentence by sentence
	for _, para := range script {
		var translated = []string{}
		for _, sentences := range para {
			//update the status entry
			entry.SetText(createProgressTemplate(currentSenetence, totalSenetences))
			line := translationclient.TranslateSentence(sentences, locale)
			translated = append(translated, line)
			currentSenetence++
		}
		translatedScript = append(translatedScript, translated)
	}

	unmarshalledScript := utils.UnmarshallScript(translatedScript)
	//set the script
	entry.SetText(unmarshalledScript)

}

func createProgressTemplate(current, total int) string {
	return fmt.Sprintf("Translating - %d/%d", current, total)
}
