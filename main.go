package main

import (
	"toy-app/dragonball"
	"toy-app/translator"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Toy App")

	dragonball.SpeakerInit()

	myWindow.SetContent(translator.GuiContent())
	myWindow.ShowAndRun()
}
