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

const padding = 5

type applet struct {
	name    string
	content func(fyne.Size) fyne.CanvasObject
	minSize fyne.Size
}

type layout struct {
	minSize fyne.Size
}

func (l *layout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return l.minSize
}

func (l *layout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	objects[0].Move(fyne.Position{
		X: 0,
		Y: 0,
	})
	objects[0].Resize(fyne.Size{
		Width:  objects[0].MinSize().Width + padding,
		Height: containerSize.Height,
	})

	objects[1].Move(fyne.Position{
		X: objects[0].Position().X + objects[0].Size().Width,
		Y: 0,
	})
	objects[1].Resize(fyne.Size{
		Width:  objects[1].MinSize().Width,
		Height: containerSize.Height,
	})

	objects[2].Move(fyne.Position{
		X: objects[1].Position().X + objects[1].Size().Width + padding,
		Y: 0,
	})
	objects[2].Resize(fyne.Size{
		Width:  containerSize.Width - objects[2].Position().X,
		Height: containerSize.Height,
	})
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Toy App")

	var applets = []applet{
		{name: "Metronome", content: metronome.Content},
		{name: "Morse translator", content: translator.Content},
		{name: "Tic-Tac-Toe", content: tictactoe.Content},
		{name: "Welcome", content: welcome.Content},
	}

	template := ""
	for idx, label, maxWidth := 0, widget.NewLabel(""), float32(0); idx < len(applets); idx += 1 {
		label.SetText(applets[idx].name)
		if width := label.MinSize().Width; width > maxWidth {
			template = applets[idx].name
			maxWidth = width
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

	separator := widget.NewSeparator()

	maxSize := fyne.Size{}
	for idx := range applets {
		applets[idx].minSize = applets[idx].content(fyne.Size{}).MinSize()
		maxSize.Width = max(maxSize.Width, applets[idx].minSize.Width)
		maxSize.Height = max(maxSize.Height, applets[idx].minSize.Height)
	}

	maxSize.Width = list.MinSize().Width + separator.MinSize().Width + maxSize.Width + 2*padding
	maxSize.Height = max(list.MinSize().Height, maxSize.Height)

	list.OnSelected = func(id widget.ListItemID) {
		myWindow.SetContent(
			container.New(&layout{maxSize}, list, separator, applets[id].content(applets[id].minSize)),
		)
	}

	list.OnUnselected = func(id widget.ListItemID) {
		speaker.Clear()
		runtime.GC()
	}

	myWindow.SetContent(
		container.New(&layout{maxSize}, list, separator, welcome.Content(applets[len(applets)-1].minSize)),
	)
	// myWindow.SetContent(applets[0].content(applets[0].minSize))
	dragonball.SpeakerInit()
	myWindow.ShowAndRun()
}
