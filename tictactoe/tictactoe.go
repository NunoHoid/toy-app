package tictactoe

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type applet struct{}

func (a applet) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

func (a applet) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {

}

type tile struct {
	widget.Icon
	row      int
	col      int
	onTapped func(idx int, jdx int)
}

func Content() fyne.CanvasObject {
	return widget.NewLabel("To be implemented!")
}
