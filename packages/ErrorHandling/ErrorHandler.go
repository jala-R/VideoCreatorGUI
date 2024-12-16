package errorhandling

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/jala-R/VideoAutomatorGUI/packages/status"
)

var ErrorBox *widget.Entry

func init() {
	ErrorBox = widget.NewMultiLineEntry()

}

func HandleErrorPop(err error) {
	handlerError(err, true)
}

func HandleError(err error) {
	handlerError(err, false)
}

func handlerError(err error, toPop bool) {
	fmt.Println(err)
	prevError := ErrorBox.Text
	errorMessage := fmt.Sprintf("%s\n%s", fmt.Errorf("%s : error occured : %w", time.Now(), err).Error(), prevError)
	ErrorBox.SetText(errorMessage)
	if toPop {
		status.ErrorPop(errorMessage)
	}
}
