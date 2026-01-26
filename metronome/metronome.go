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

type applet struct{}

func (a applet) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(500, 200)
}

func (a applet) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	const offset = 20

	objects[0].Move(fyne.NewPos(
		(containerSize.Width-a.MinSize(nil).Width)/2,
		(containerSize.Height-a.MinSize(nil).Height)/2,
	))
	objects[0].Resize(a.MinSize(nil))

	objects[1].Move(fyne.NewPos(objects[0].Position().X+offset, containerSize.Height/2-3*offset))
	objects[1].Resize(fyne.NewSquareSize(2 * offset))

	objects[2].Move(objects[1].Position().AddXY(0, 4*offset))
	objects[2].Resize(fyne.NewSquareSize(2 * offset))

	objects[3].Move(objects[0].Position().AddXY(4*offset, offset))
	objects[4].Resize(objects[4].MinSize())

	objects[3].Resize(a.MinSize(nil).SubtractWidthHeight(objects[4].Size().Width+6*offset, 2*offset))
	objects[4].Move(fyne.NewPos(
		objects[3].Position().X+objects[3].Size().Width+offset,
		(containerSize.Height-objects[4].Size().Height)/2,
	))
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

func Content() fyne.CanvasObject {
	beatsSelect := widget.NewSelect([]string{"2", "3", "4", "6"}, nil)
	beatsSelect.SetSelected("4")

	borderCard := widget.NewCard("", "", nil)

	speedCard := widget.NewCard("", "", widget.NewSlider(60, 300))
	speedCard.Content.(*widget.Slider).OnChanged = func(f float64) {
		speedCard.SetTitle(fmt.Sprint(f))
	}
	speedCard.Content.(*widget.Slider).SetValue(120)
	speedCard.Content.(*widget.Slider).Step = 5

	beatsRadio := widget.NewRadioGroup([]string{"2", "3", "4"}, nil)
	beatsRadio.SetSelected("4")

	playButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		speaker.Clear()
		playOneHour(int(speedCard.Content.(*widget.Slider).Value), int(beatsRadio.Selected[0]-'0'))
	})

	stopButton := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
		speaker.Clear()
	})

	return container.New(applet{}, borderCard, playButton, stopButton, speedCard, beatsRadio)
}
