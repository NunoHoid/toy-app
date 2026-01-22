package welcome

import (
	_ "embed"
	"toy-app/dragonball"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//go:embed mascot.png
var mascot []byte

type applet struct {
	image *canvas.Image
}

func (a applet) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return objects[0].MinSize().AddWidthHeight(0, float32(a.image.Image.Bounds().Max.Y))
}

func (a applet) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	objects[0].Move(fyne.NewPos(
		(containerSize.Width-objects[0].Size().Width)/2,
		(containerSize.Height-objects[0].Size().Height)/2,
	))
	objects[0].Resize(a.MinSize(objects))
}

func GuiContent() fyne.CanvasObject {
	image := canvas.NewImageFromReader(dragonball.NewFile(mascot), "mascot.png")
	image.Resize(fyne.NewSquareSize(1))
	return container.New(applet{image}, widget.NewCard(
		"Welcome",
		"A toy app to explore the fyne framework",
		image,
	))
}
