package main

import (
	"runtime"
	"toy-app/dragonball"
	"toy-app/metronome"
	"toy-app/tictactoe"
	"toy-app/translator"
	"toy-app/welcome"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gopxl/beep/v2/speaker"
)

var applets = []applet{
	{name: "Metronome", content: metronome.GuiContent},
	{name: "Morse translator", content: translator.GuiContent},
	{name: "Tic-Tac-Toe", content: tictactoe.GuiContent},
}

var template = ""

type applet struct {
	name    string
	content func() fyne.CanvasObject
}

func (_ applet) MinSize(objects []fyne.CanvasObject) fyne.Size {
	maxSize := fyne.NewSize(0, 0)
	for _, val := range applets {
		minSize := val.content().MinSize()
		maxSize.Width = max(maxSize.Width, minSize.Width)
		maxSize.Height = max(maxSize.Height, minSize.Height)
	}
	return maxSize.AddWidthHeight(widget.NewLabel(template).MinSize().Width, 0)
}

func (_ applet) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(fyne.NewSize(objects[0].MinSize().Width+5, containerSize.Height))

	objects[1].Move(fyne.NewPos(objects[0].Size().Width, 0))
	objects[1].Resize(fyne.NewSize(objects[1].MinSize().Width, containerSize.Height))

	objects[2].Move(fyne.NewPos(objects[0].Size().Width+objects[1].Size().Width+5, 0))
	objects[2].Resize(fyne.NewSize(containerSize.Width-objects[2].Position().X, containerSize.Height))
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Toy App")

	for idx, val := 0, float32(0); idx < len(applets); idx += 1 {
		if width := widget.NewLabel(applets[idx].name).MinSize().Width; width > val {
			template = applets[idx].name
			val = width
		}
	}

	list := widget.NewList(
		func() int {
			return len(applets)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(template)
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(applets[id].name)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		speaker.Clear()
		myWindow.SetContent(container.New(applet{}, list, widget.NewSeparator(), applets[id].content()))
		runtime.GC()
	}

	myWindow.SetContent(container.New(applet{}, list, widget.NewSeparator(), welcome.GuiContent()))
	dragonball.SpeakerInit()
	myWindow.ShowAndRun()
}
