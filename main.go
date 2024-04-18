package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Input struct {
	widget.Entry
	Value string
}

func NewInput() *Input {
	e := &Input{}
	e.ExtendBaseWidget(e)
	e.OnChanged = func(text string) {
		e.Value = text
	}
	return e
}

func main() {
	app := app.New()
	window := app.NewWindow("App Installer")
	window.Resize(fyne.Size{Width: 400, Height: 300})
	app.Settings().SetTheme(theme.DarkTheme())

	var fullFilePath string
	var fileName string
	var userInputAppName string
	var fullIconPath string
	var fileIconName string

	var username string
	var dirPath string
	var fileBaseName string
	var filePath string
	// var parts []string

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Add name to your app *(Ignore for setting existing name)")

	descriptionLabelFile := widget.NewLabel("Choose app: ")
	isSpecifiedAppLabel := widget.NewLabel("")

	captureInput := func() {
		userInputAppName = entry.Text
		if !(len(fileName) < 1) {
			isSpecifiedAppLabel.SetText("*Saved")
		} else {
			isSpecifiedAppLabel.SetText("*Nothing to save")
		}
	}

	captureButton := widget.NewButton("Save", captureInput)

	openFileButton := widget.NewButton("Choose file you want to add to overview", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if r == nil {
				return
			}

			fullFilePath = r.URI().Path()
			fileName = filepath.Base(fullFilePath)

			fmt.Println("File path: ", fullFilePath)

			descriptionLabelFile.SetText("Chosen app: " + fileName)

		}, window)
	})

	descriptionLabelIcon := widget.NewLabel("Choose icon: *(Ignore unless you want fully transparent one) ")

	openIconButton := widget.NewButton("Choose icon", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if r == nil {
				return
			}

			fullIconPath = r.URI().Path()
			fileIconName = filepath.Base(fullIconPath)

			fmt.Println("Icon path: ", fullFilePath)

			descriptionLabelIcon.SetText("Chosen icon: " + fileIconName)
		}, window)
	})

	user_password := NewInput()

	items := []*widget.FormItem{
		widget.NewFormItem("Password", user_password),
	}

	createDotConfigButton := widget.NewButton("Create app", func() {
		if len(fullFilePath) < 1 {
			isSpecifiedAppLabel.SetText("*You didn't specify the path")
			return
		}
		if len(userInputAppName) < 1 {
			userInputAppName = fileName
		}
		config := fmt.Sprintf(`#!/usr/bin/env xdg-open
[Desktop Entry]
Version=1.0
Type=Application
Terminal=false
Exec=%s
Name=%s
Comment=
Icon=%s`, fullFilePath, userInputAppName, fullIconPath)

		currentUser, err := user.Current()
		if err != nil {
			log.Fatalf(err.Error())
		}
		username = currentUser.Username

		dirPath = fmt.Sprintf("/home/%s/Documents", username)

		fileBaseName = fmt.Sprintf("%s.desktop", fileName)

		filePath = dirPath + "/" + fileBaseName

		file, err := os.Create(filePath)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		defer file.Close()

		_, err = file.WriteString(config)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowForm("Enter your password...", "Enter", "Cancel", items, func(_ bool) {

			commandToExecute := fmt.Sprintf("sudo mv %s /usr/share/applications/ && sudo chmod +x /usr/share/applications/%s.desktop", filePath, fileName)

			cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | sudo -S %s ", user_password.Value, commandToExecute))

			output, err := cmd.Output()
			isSpecifiedAppLabel.SetText("*App has been successfully added")

			dialog.ShowInformation("info", string(output), window)
			if err != nil {

				return
			}
		}, window)

	})

	window.SetContent(
		container.NewVBox(
			descriptionLabelFile,
			openFileButton,
			descriptionLabelIcon,
			openIconButton,
			entry,
			captureButton,
			isSpecifiedAppLabel,
			container.NewBorder(layout.NewSpacer(), createDotConfigButton, layout.NewSpacer(), layout.NewSpacer()),
		),
	)

	window.ShowAndRun()
}
