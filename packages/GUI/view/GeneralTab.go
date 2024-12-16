package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/controller"
	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/view/types"
)

//general
//create projects
//set project name done
//set output folder done
//set images folder done
//select script
//select languages
//uploaded file -> //english
//german
//spanish
//portugal
//convert audio
//

//general
//create projects
//set project name done
//set output folder done
//set images folder done
//set audio folder

func generalTabGUI(w fyne.Window) fyne.CanvasObject {

	errorhandling.ErrorBox.SetMinRowsVisible(20)
	generalScreenForm := widget.NewForm(
		widget.NewFormItem("Project Name", createSingleEntry("Enter project Name", controller.ProjectNameHandler)),
		widget.NewFormItem("Output Folder", selectResourceDialog(types.NewFolderType("Browse", w, controller.ProjectOutputFolder))),
		widget.NewFormItem("Image Folder", selectResourceDialog(types.NewFolderType("Browse", w, controller.ImagesFolder))),
		widget.NewFormItem("Audio Folder", selectResourceDialog(types.NewFolderType("Browse", w, controller.ReuseAudioFolder))),
		widget.NewFormItem("Images order", widget.NewCheck("Inorder", controller.ImagesInOrder)),
		widget.NewFormItem("Sentence Gap", createSingleEntry("0.0", controller.SentenceGap)),
		widget.NewFormItem("Para Gap", createSingleEntry("0.0", controller.ParaGap)),

		widget.NewFormItem("", errorhandling.ErrorBox),
	)

	generalScreenForm.OnSubmit = controller.MainScreenSubmit(w)

	return generalScreenForm
}
