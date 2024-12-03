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
	ErrorBox.SetText(fmt.Errorf("error occured : %w", err).Error())
}
