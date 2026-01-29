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

type layout struct {
	minSize fyne.Size
}

func (l *layout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if !l.minSize.IsZero() {
		return l.minSize
	}

	image := objects[0].(*widget.Card).Content.(*canvas.Image).Image

	return fyne.Size{
		Width:  max(objects[0].MinSize().Width, float32(image.Bounds().Max.X)),
		Height: objects[0].MinSize().Height + float32(image.Bounds().Max.Y),
	}
}

func (l *layout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	objects[0].Move(fyne.Position{
		X: (containerSize.Width - l.minSize.Width) / 2,
		Y: (containerSize.Height - l.minSize.Height) / 2,
	})
	objects[0].Resize(fyne.Size{
		Width:  l.minSize.Width,
		Height: l.minSize.Height,
	})
}

func Content(minSize fyne.Size) fyne.CanvasObject {
	image, _ := png.Decode(dragonball.NewFile(mascot))
	card := widget.NewCard("Welcome", "A toy app to explore the fyne framework", canvas.NewImageFromImage(image))
	card.Content.(*canvas.Image).FillMode = canvas.ImageFillContain

	return container.New(&layout{minSize}, card)
}
