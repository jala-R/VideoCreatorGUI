package controller

import (
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	apiclient "github.com/jala-R/VideoAutomatorGUI/packages/ApiClient"
	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
)

func init() {
	model.AddToDb(model.IMAGESINORDER, false)
}

func ProjectNameHandler(projName string) {
	model.AddToDb(model.PROJECTNAME, projName)
}

func ProjectOutputFolder(outFolder string) {
	model.AddToDb(model.OUTPUTFOLDER, outFolder)
}

func ImagesFolder(imgFolder string) {
	model.AddToDb(model.IMAGEFOLDER, imgFolder)
}

func ReuseAudioFolder(folder string) {
	model.AddToDb(model.REUSEAUDIOFOLDER, folder)
}

func ImagesInOrder(state bool) {
	model.AddToDb(model.IMAGESINORDER, state)
}

func SentenceGap(gap string) {
	gapInFloat, err := strconv.ParseFloat(gap, 64)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}
	model.AddToDb(model.SENTENCEGAP, gapInFloat)
}

func ParaGap(gap string) {
	gapInFloat, err := strconv.ParseFloat(gap, 64)
	if err != nil {
		errorhandling.HandleErrorPop(err)
		return
	}
	model.AddToDb(model.PARAGAP, gapInFloat)
}

func ValidateGeneralSceenFeilds() error {
	if IsEmpty(model.QueryDB(model.PROJECTNAME)) {
		return errors.New("enter project  name")
	}
	if IsEmpty(model.QueryDB(model.OUTPUTFOLDER)) {
		return errors.New("select output folder")
	}
	if IsEmpty(model.QueryDB(model.IMAGEFOLDER)) {
		return errors.New("select image folder ")
	}
	if IsEmpty(model.QueryDB(model.REUSEAUDIOFOLDER)) {
		return errors.New("select audio folder")
	}

	return nil
}

func MainScreenSubmit(w fyne.Window) func() {
	return func() {
		fmt.Println("Main screen submitted")

		err := ValidateGeneralSceenFeilds()
		if err != nil {
			errorhandling.HandleErrorPop(err)
			return
		}

		apiclient.CreateVideoProject()

		info := dialog.NewInformation("Status", "Project Completed", w)
		info.Resize(fyne.NewSize(200, 200))
		info.Show()
	}
}
