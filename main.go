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
	{name: "Metronome", content: metronome.Content},
	{name: "Morse translator", content: translator.Content},
	{name: "Tic-Tac-Toe", content: tictactoe.Content},
	{name: "Welcome", content: welcome.Content},
}

const offset = 5

type applet struct {
	name    string
	content func() fyne.CanvasObject
}

func (a applet) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := objects[0].MinSize().Add(objects[1].MinSize()).AddWidthHeight(0, 2*offset)
	maxSize := fyne.NewSquareSize(0)
	for _, val := range applets {
		maxSize = maxSize.Max(val.content().MinSize())
	}
	return minSize.Add(maxSize)
}

func (a applet) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(fyne.NewSize(objects[0].MinSize().Width+offset, containerSize.Height))

	objects[1].Move(fyne.NewPos(objects[0].Size().Width, 0))
	objects[1].Resize(fyne.NewSize(objects[1].MinSize().Width, containerSize.Height))

	objects[2].Move(fyne.NewPos(objects[1].Position().X+objects[1].Size().Width+offset, 0))
	objects[2].Resize(fyne.NewSize(containerSize.Width-objects[2].Position().X, containerSize.Height))
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Toy App")

	template := ""
	for idx, val := 0, float32(0); idx < len(applets); idx += 1 {
		if width := widget.NewLabel(applets[idx].name).MinSize().Width; width > val {
			template = applets[idx].name
			val = width
		}
	}

	list := widget.NewList(
		func() int {
			return len(applets) - 1
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(template)
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(applets[id].name)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		myWindow.SetContent(container.New(applet{}, list, widget.NewSeparator(), applets[id].content()))
	}

	list.OnUnselected = func(id widget.ListItemID) {
		speaker.Clear()
		runtime.GC()
	}

	myWindow.SetContent(container.New(applet{}, list, widget.NewSeparator(), welcome.Content()))
	dragonball.SpeakerInit()
	myWindow.ShowAndRun()
}
