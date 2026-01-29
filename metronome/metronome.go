package metronome

import (
	"fmt"
	"strconv"
	"toy-app/dragonball"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gopxl/beep/v2/midi"
	"github.com/gopxl/beep/v2/speaker"
)

const padding = 5

type layout struct {
	minSize fyne.Size
}

func (l *layout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if !l.minSize.IsZero() {
		return l.minSize
	}

	return fyne.Size{
		Width:  2 * objects[1].MinSize().Width,
		Height: objects[0].MinSize().Height + objects[1].MinSize().Height + padding,
	}
}

func (l *layout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	objects[0].Move(fyne.Position{
		X: (containerSize.Width - l.minSize.Width) / 2,
		Y: (containerSize.Height - l.minSize.Height) / 2,
	})
	objects[0].Resize(fyne.Size{
		Width:  l.minSize.Width,
		Height: objects[0].MinSize().Height,
	})

	objects[1].Move(fyne.Position{
		X: objects[0].Position().X,
		Y: objects[0].Position().Y + objects[0].Size().Height + padding,
	})
	objects[1].Resize(fyne.Size{
		Width:  objects[1].MinSize().Width,
		Height: objects[1].MinSize().Height,
	})

	objects[2].Move(fyne.Position{
		X: objects[1].Position().X + 6*objects[1].Size().Width/5,
		Y: objects[1].Position().Y + (objects[1].Size().Height-objects[1].Size().Width/5)/2,
	})
	objects[2].Resize(fyne.Size{
		Width:  objects[1].Size().Width / 5,
		Height: objects[1].Size().Width / 5,
	})

	objects[3].Move(fyne.Position{
		X: objects[2].Position().X + 2*objects[1].Size().Width/5,
		Y: objects[2].Position().Y,
	})
	objects[3].Resize(fyne.Size{
		Width:  objects[1].Size().Width / 5,
		Height: objects[1].Size().Width / 5,
	})
}

func playOneHour(beatsPerMinute int, beatsPerBar int) {
	content := []byte{
		0x4d, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x01,
		0x00, 0x01,
		0x01, 0xe0,
		0x4d, 0x54, 0x72, 0x6b,
		0x00, 0x00, 0x00, 0x21,
		0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20,
		0x00, 0xb0, 0x07, 0xff,
		0x00, 0xc0, 0x7f,
		0x00, 0x90, 0x18, 0x50, 0x83, 0x5f, 0x18,
	}

	speed := fmt.Sprintf("%06x", 60_000_000/beatsPerMinute)
	for idx := range 3 {
		val, _ := strconv.ParseInt(speed[2*idx:2*idx+2], 16, 0)
		content[idx+26] = byte(val)
	}

	for idx, beat := 1, 1; idx < 60*beatsPerMinute; idx, beat = idx+1, beat+1 {
		if beat == beatsPerBar {
			content = append(content, 0x00, 0x01, 0x18, 0x50, 0x83, 0x5f, 0x18)
			beat = 0
		} else {
			content = append(content, 0x00, 0x01, 0x2c, 0x50, 0x83, 0x5f, 0x2c)
		}
	}

	content = append(content, 0x00, 0x01, 0xff, 0x2f, 0x00)

	stream, _, _ := midi.Decode(dragonball.NewFile(content), dragonball.NewFont(), dragonball.SampleRate)
	speaker.Play(stream)
}

func Content(minSize fyne.Size) fyne.CanvasObject {
	speedCard := widget.NewCard("", "Beats per minute", widget.NewSlider(60, 300))
	speedCard.Content.(*widget.Slider).OnChanged = func(f float64) {
		speedCard.SetTitle(fmt.Sprint(f))
	}
	speedCard.Content.(*widget.Slider).SetValue(120)
	speedCard.Content.(*widget.Slider).Step = 5

	beatsCard := widget.NewCard("", "Beats per measure", widget.NewRadioGroup([]string{"2", "3", "4"}, nil))
	beatsCard.Content.(*widget.RadioGroup).SetSelected("4")
	beatsCard.Content.(*widget.RadioGroup).Horizontal = true

	playButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		speaker.Clear()
		playOneHour(
			int(speedCard.Content.(*widget.Slider).Value),
			int(beatsCard.Content.(*widget.RadioGroup).Selected[0]-'0'),
		)
	})

	stopButton := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
		speaker.Clear()
	})

	return container.New(&layout{minSize}, speedCard, beatsCard, playButton, stopButton)
}
