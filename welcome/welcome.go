package welcome

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func GuiContent() fyne.CanvasObject {
	return widget.NewLabel("Welcome!")
}
