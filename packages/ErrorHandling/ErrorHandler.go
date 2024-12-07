package errorhandling

import (
	"fmt"

	"fyne.io/fyne/v2/widget"
)

var ErrorBox *widget.Entry

func init() {
	ErrorBox = widget.NewMultiLineEntry()

}

func HandleError(err error) {
	fmt.Println(err)
	prevError := ErrorBox.Text
	errorMessage := fmt.Sprintf("%s\n%s", fmt.Errorf("error occured : %w", err).Error(), prevError)
	ErrorBox.SetText(errorMessage)
}
