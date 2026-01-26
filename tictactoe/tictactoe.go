package tictactoe

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var state = []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
var winnerLabel = widget.NewLabel("Circle plays")

type applet struct{}

func (a applet) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(400, 400)
}

func (a applet) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	const offset = 40

	tileSize := (min(containerSize.Width, containerSize.Height) - offset) / 3

	objects[9].Move(fyne.NewPos((containerSize.Width-3*tileSize-offset)/2, (containerSize.Height-3*tileSize-offset)/2))
	objects[10].Move(objects[9].Position().AddXY(offset, 0))

	objects[9].Resize(fyne.NewSquareSize(offset))
	objects[10].Resize(fyne.NewSize(3*tileSize, offset))

	objects[0].Move(objects[9].Position().Add(objects[9].Size()))
	objects[1].Move(objects[0].Position().AddXY(tileSize, 0))
	objects[2].Move(objects[1].Position().AddXY(tileSize, 0))

	objects[3].Move(objects[0].Position().AddXY(0, tileSize))
	objects[4].Move(objects[3].Position().AddXY(tileSize, 0))
	objects[5].Move(objects[4].Position().AddXY(tileSize, 0))

	objects[6].Move(objects[3].Position().AddXY(0, tileSize))
	objects[7].Move(objects[6].Position().AddXY(tileSize, 0))
	objects[8].Move(objects[7].Position().AddXY(tileSize, 0))

	for _, val := range objects[:9] {
		val.Resize(fyne.NewSquareSize(tileSize))
	}
}

type tile struct {
	widget.Card
	idx    int
	jdx    int
	isDown bool
}

func newTile(idx int, jdx int) fyne.CanvasObject {
	card := &tile{idx: idx, jdx: jdx, isDown: false}
	card.ExtendBaseWidget(card)
	return card
}

func (t *tile) MouseIn(event *desktop.MouseEvent) {
}

func (t *tile) MouseOut() {
	t.isDown = false
}

func (t *tile) MouseMoved(event *desktop.MouseEvent) {
}

func (t *tile) MouseDown(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonPrimary {
		t.isDown = true
	}
}

func (t *tile) MouseUp(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonPrimary && state[0] < 9 && t.isDown && t.Content == nil {
		if state[0]%2 == 0 {
			t.SetContent(canvas.NewImageFromResource(theme.RadioButtonIcon()))
			winnerLabel.SetText("Cross plays")
		} else {
			t.SetContent(canvas.NewImageFromResource(theme.CancelIcon()))
			winnerLabel.SetText("Circle plays")
		}
		state[0] += 1
		state[t.idx+1] += 2*(state[0]%2) - 1
		state[t.jdx+4] += 2*(state[0]%2) - 1
		if t.idx == t.jdx {
			state[7] += 2*(state[0]%2) - 1
		}
		if t.idx+t.jdx == 2 {
			state[8] += 2*(state[0]%2) - 1
		}
		for _, val := range state[1:] {
			if val == 3 {
				winnerLabel.SetText("Circle wins")
				state[0] = 9
				return
			}
			if val == -3 {
				winnerLabel.SetText("Cross wins")
				state[0] = 9
				return
			}
		}
		if state[0] == 9 {
			winnerLabel.SetText("Tie")
		}
	}
	t.isDown = false
}

func Content() fyne.CanvasObject {
	winnerLabel.Alignment = fyne.TextAlignCenter

	tiles := []fyne.CanvasObject{
		newTile(0, 0), newTile(0, 1), newTile(0, 2),
		newTile(1, 0), newTile(1, 1), newTile(1, 2),
		newTile(2, 0), newTile(2, 1), newTile(2, 2),
	}

	replayButton := widget.NewButtonWithIcon("", theme.MediaReplayIcon(), func() {
		clear(state)
		winnerLabel.SetText("Circle plays")
		for _, val := range tiles {
			val.(*tile).SetContent(nil)
		}
	})

	return container.New(applet{}, append(tiles, replayButton, winnerLabel)...)
}
