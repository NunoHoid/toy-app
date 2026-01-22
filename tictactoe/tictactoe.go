package tictactoe

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type tile struct {
	widget.Icon
	row      int
	col      int
	onTapped func(idx int, jdx int)
}

func GuiContent() fyne.CanvasObject {
	return widget.NewLabel("To be implemented!")
}
