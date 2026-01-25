package welcome

import (
	_ "embed"
	"image/png"
	"toy-app/dragonball"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//go:embed mascot.png
var mascot []byte

type applet struct{}

func (a applet) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(
		max(objects[0].MinSize().Width, float32(objects[0].(*widget.Card).Content.(*canvas.Image).Image.Bounds().Max.X)),
		objects[0].MinSize().Height+float32(objects[0].(*widget.Card).Content.(*canvas.Image).Image.Bounds().Max.Y),
	)
}

func (a applet) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	objects[0].Move(fyne.NewPos(
		(containerSize.Width-a.MinSize(objects).Width)/2,
		(containerSize.Height-a.MinSize(objects).Height)/2,
	))
	objects[0].Resize(a.MinSize(objects))
}

func Content() fyne.CanvasObject {
	image, _ := png.Decode(dragonball.NewFile(mascot))
	card := widget.NewCard("Welcome", "A toy app to explore the fyne framework", canvas.NewImageFromImage(image))
	card.Content.(*canvas.Image).FillMode = canvas.ImageFillContain
	return container.New(applet{}, card)
}
