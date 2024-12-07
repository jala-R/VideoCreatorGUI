package types

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
)

type ResourceType interface {
	CreateButton(label *widget.Label) *widget.Button
	Flush(label *widget.Label)
}

type FolderType struct {
	name    string
	w       fyne.Window
	trigger func(string)
}

func NewFolderType(name string, w fyne.Window, trigger func(string)) *FolderType {
	return &FolderType{
		name:    name,
		w:       w,
		trigger: trigger,
	}
}

func (folder *FolderType) CreateButton(label *widget.Label) *widget.Button {
	button := widget.NewButton(folder.name, func() {
		fd := dialog.NewFolderOpen(folder.onSelect(label), folder.w)

		fd.Resize(fyne.NewSize(1000, 1000))
		fd.Show()
	})

	return button
}

func (folder *FolderType) Flush(label *widget.Label) {
	label.SetText("")
	label.Hide()
	folder.trigger("")
}

func (folder *FolderType) onSelect(label *widget.Label) func(fyne.ListableURI, error) {
	return func(url fyne.ListableURI, err error) {
		if err != nil {
			errorhandling.HandleError(err)
			return
		}

		if url != nil {
			label.SetText(url.Path())
			label.Show()

			folder.trigger(url.Path())
		}

	}
}

type FileType struct {
	FolderType
	ext string
}

func NewFileType(name string, w fyne.Window, trigger func(string), ext string) *FileType {
	return &FileType{
		FolderType: *NewFolderType(name, w, trigger),
		ext:        ext,
	}
}

func (file *FileType) CreateButton(label *widget.Label) *widget.Button {
	button := widget.NewButton(file.name, func() {
		fd := dialog.NewFileOpen(file.onSelect(label), file.w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{file.ext}))
		fd.Resize(fyne.NewSize(1000, 1000))
		fd.Show()
	})

	return button
}

func (file *FileType) onSelect(label *widget.Label) func(fyne.URIReadCloser, error) {
	return func(url fyne.URIReadCloser, err error) {

		if err != nil {
			errorhandling.HandleError(err)
			return
		}
		if url == nil {
			errorhandling.HandleError(errors.New("no file selected"))
			return
		}

		defer url.Close()

		path := url.URI().Path()

		label.SetText(string(path))
		label.Show()

		file.trigger(string(path))

	}
}
