package tictactoe

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type layout struct {
	minSize fyne.Size
}

func (l *layout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if !l.minSize.IsZero() {
		return l.minSize
	}

	return fyne.Size{
		Width:  objects[10].MinSize().Width + objects[10].MinSize().Height,
		Height: objects[10].MinSize().Width + objects[10].MinSize().Height,
	}
}

func (l *layout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	topHeight := objects[10].MinSize().Height

	tileSize := fyne.Size{
		Width:  (min(containerSize.Width, containerSize.Height) - topHeight) / 3,
		Height: (min(containerSize.Width, containerSize.Height) - topHeight) / 3,
	}

	objects[9].Move(fyne.Position{
		X: (containerSize.Width - 3*tileSize.Width - topHeight) / 2,
		Y: (containerSize.Height - 3*tileSize.Height - topHeight) / 2,
	})
	objects[9].Resize(fyne.Size{
		Width:  topHeight,
		Height: topHeight,
	})

	objects[10].Move(fyne.Position{
		X: objects[9].Position().X + topHeight,
		Y: objects[9].Position().Y,
	})
	objects[10].Resize(fyne.Size{
		Width:  3 * tileSize.Width,
		Height: topHeight,
	})

	for idx, val := range objects[:9] {
		val.Move(fyne.Position{
			X: objects[9].Position().X + float32(idx%3)*tileSize.Width + topHeight,
			Y: objects[9].Position().Y + float32(idx/3)*tileSize.Height + topHeight,
		})
		val.Resize(fyne.Size{
			Width:  tileSize.Width,
			Height: tileSize.Height,
		})
	}
}

type tile struct {
	widget.Card
	idx    int
	jdx    int
	isDown bool
	score  []int
	label  *widget.Label
}

func (t *tile) MouseIn(event *desktop.MouseEvent) {
}

func (t *tile) MouseOut() {
	t.isDown = false
}

func (t *tile) MouseMoved(event *desktop.MouseEvent) {
}

func (t *tile) MouseDown(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonPrimary && t.Content == nil && t.score[0] < 9 {
		t.isDown = true
	}
}

func (t *tile) MouseUp(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonPrimary && t.isDown {
		if t.score[0]%2 == 0 {
			t.SetContent(canvas.NewImageFromResource(theme.RadioButtonIcon()))
			t.label.SetText("Cross plays")
		} else {
			t.SetContent(canvas.NewImageFromResource(theme.CancelIcon()))
			t.label.SetText("Circle plays")
		}
		t.score[0] += 1
		t.score[t.idx+1] += 2*(t.score[0]%2) - 1
		t.score[t.jdx+4] += 2*(t.score[0]%2) - 1
		if t.idx == t.jdx {
			t.score[7] += 2*(t.score[0]%2) - 1
		}
		if t.idx+t.jdx == 2 {
			t.score[8] += 2*(t.score[0]%2) - 1
		}
		for _, val := range t.score[1:] {
			if val == 3 {
				t.label.SetText("Circle wins")
				t.score[0] = 9
				return
			}
			if val == -3 {
				t.label.SetText("Cross wins")
				t.score[0] = 9
				return
			}
		}
		if t.score[0] == 9 {
			t.label.SetText("It's a tie")
		}
	}
	t.isDown = false
}

func newTile(idx int, jdx int, score []int, label *widget.Label) fyne.CanvasObject {
	card := &tile{idx: idx, jdx: jdx, isDown: false, score: score, label: label}
	card.ExtendBaseWidget(card)
	return card
}

func Content(minSize fyne.Size) fyne.CanvasObject {
	score := [9]int{}
	topLabel := widget.NewLabel("Circle plays")

	tiles := []fyne.CanvasObject{
		newTile(0, 0, score[:], topLabel), newTile(0, 1, score[:], topLabel), newTile(0, 2, score[:], topLabel),
		newTile(1, 0, score[:], topLabel), newTile(1, 1, score[:], topLabel), newTile(1, 2, score[:], topLabel),
		newTile(2, 0, score[:], topLabel), newTile(2, 1, score[:], topLabel), newTile(2, 2, score[:], topLabel),
	}

	replayButton := widget.NewButtonWithIcon("", theme.MediaReplayIcon(), func() {
		clear(score[:])
		topLabel.SetText("Circle plays")
		for _, val := range tiles {
			val.(*tile).SetContent(nil)
		}
	})

	topLabel.Alignment = fyne.TextAlignCenter

	return container.New(&layout{minSize}, append(tiles, replayButton, topLabel)...)
}
