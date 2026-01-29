package translator

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"toy-app/dragonball"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gopxl/beep/v2/midi"
	"github.com/gopxl/beep/v2/speaker"
)

const padding = 20

type layout struct {
	minSize fyne.Size
}

func (l *layout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if !l.minSize.IsZero() {
		return l.minSize
	}

	centerHeight := float32(0)
	for _, val := range objects {
		centerHeight = max(centerHeight, val.MinSize().Height)
	}

	centerWidth := objects[2].MinSize().Width + objects[5].MinSize().Width + objects[6].MinSize().Width
	centerWidth += 3*centerHeight + 6*padding

	return fyne.Size{
		Width:  max(centerWidth, objects[0].MinSize().Width),
		Height: 3*centerHeight + 2*padding,
	}
}

func (l *layout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	centerSize := fyne.Size{
		Width:  l.minSize.Width - 2*padding,
		Height: (l.minSize.Height-2*padding)/3 - 2*padding,
	}

	entrySize := fyne.Size{
		Width:  containerSize.Width,
		Height: (containerSize.Height-centerSize.Height)/2 - padding,
	}

	objects[0].Move(fyne.Position{
		X: 0,
		Y: 0,
	})
	objects[0].Resize(fyne.Size{
		Width:  entrySize.Width,
		Height: entrySize.Height,
	})

	objects[1].Move(fyne.Position{
		X: 0,
		Y: containerSize.Height - entrySize.Height,
	})
	objects[1].Resize(fyne.Size{
		Width:  entrySize.Width,
		Height: entrySize.Height,
	})

	objects[2].Move(fyne.Position{
		X: (containerSize.Width - centerSize.Width) / 2,
		Y: entrySize.Height + padding,
	})
	objects[2].Resize(fyne.Size{
		Width:  objects[2].MinSize().Width,
		Height: centerSize.Height,
	})

	objects[3].Move(fyne.Position{
		X: objects[2].Position().X + objects[2].Size().Width + padding,
		Y: objects[2].Position().Y,
	})
	objects[3].Resize(fyne.Size{
		Width:  centerSize.Height,
		Height: centerSize.Height,
	})

	objects[4].Move(fyne.Position{
		X: objects[3].Position().X + objects[3].Size().Width + padding,
		Y: objects[3].Position().Y,
	})
	objects[4].Resize(fyne.Size{
		Width:  centerSize.Height,
		Height: centerSize.Height,
	})

	objects[5].Move(fyne.Position{
		X: objects[4].Position().X + objects[4].Size().Width + padding,
		Y: objects[4].Position().Y,
	})
	objects[5].Resize(fyne.Size{
		Width:  objects[5].MinSize().Width,
		Height: centerSize.Height,
	})

	objects[6].Move(fyne.Position{
		X: objects[5].Position().X + objects[5].Size().Width + padding,
		Y: objects[5].Position().Y,
	})
	objects[6].Resize(fyne.Size{
		Width:  objects[6].MinSize().Width,
		Height: centerSize.Height,
	})

	objects[7].Move(fyne.Position{
		X: objects[6].Position().X + objects[6].Size().Width,
		Y: objects[6].Position().Y,
	})
	objects[7].Resize(fyne.Size{
		Width:  (containerSize.Width+centerSize.Width)/2 - objects[7].Position().X,
		Height: centerSize.Height,
	})
}

func letters() [][2]string {
	return [][2]string{
		{"a", ".-"},
		{"b", "-..."},
		{"c", "-.-."},
		{"d", "-.."},
		{"e", "."},
		{"f", "..-."},
		{"g", "--."},
		{"h", "...."},
		{"i", ".."},
		{"j", ".---"},
		{"k", "-.-"},
		{"l", ".-.."},
		{"m", "--"},
		{"n", "-."},
		{"o", "---"},
		{"p", ".--."},
		{"q", "--.-"},
		{"r", ".-."},
		{"s", "..."},
		{"t", "-"},
		{"u", "..-"},
		{"v", "...-"},
		{"w", ".--"},
		{"x", "-..-"},
		{"y", "-.--"},
		{"z", "--.."},
	}
}

func numbers() [][2]string {
	return [][2]string{
		{"0", "-----"},
		{"1", ".----"},
		{"2", "..---"},
		{"3", "...--"},
		{"4", "....-"},
		{"5", "....."},
		{"6", "-...."},
		{"7", "--..."},
		{"8", "---.."},
		{"9", "----."},
	}
}

func markers() [][2]string {
	return [][2]string{
		{"!", "-.-.--"},
		{"\"", ".-..-."},
		{"#", "#"},
		{"$", "#"},
		{"%", "#"},
		{"&", ".-..."},
		{"'", ".----."},
		{"(", "-.--."},
		{")", "-.--.-"},
		{"*", "#"},
		{"+", ".-.-."},
		{",", "--..--"},
		{"-", "-....-"},
		{".", ".-.-.-"},
		{"/", "-..-."},
		{":", "---..."},
		{";", "#"},
		{"<", "#"},
		{"=", "-...-"},
		{">", "#"},
		{"?", "..--.."},
		{"@", ".--.-."},
	}
}

func encodeChar(char rune) string {
	if 'A' <= char && char <= 'Z' {
		return letters()[char-'A'][1]
	}
	if 'a' <= char && char <= 'z' {
		return letters()[char-'a'][1]
	}
	if '0' <= char && char <= '9' {
		return numbers()[char-'0'][1]
	}
	if '!' <= char && char <= '/' {
		return markers()[char-'!'][1]
	}
	if ':' <= char && char <= '@' {
		return markers()[char-':'][1]
	}
	return "#"
}

func decodeChar(char string) rune {
	if idx := slices.IndexFunc(letters(), func(pair [2]string) bool { return char == pair[1] }); idx != -1 {
		return rune(letters()[idx][0][0])
	}
	if idx := slices.IndexFunc(numbers(), func(pair [2]string) bool { return char == pair[1] }); idx != -1 {
		return rune(numbers()[idx][0][0])
	}
	if idx := slices.IndexFunc(markers(), func(pair [2]string) bool { return char == pair[1] }); idx != -1 {
		return rune(markers()[idx][0][0])
	}
	return '#'
}

func clearLatin(text string) string {
	builder := strings.Builder{}
	lastIsSpace := false

	for word := range strings.FieldsSeq(text) {
		if builder.Len() > 0 && !lastIsSpace {
			builder.WriteRune(' ')
			lastIsSpace = true
		}
		for _, char := range word {
			if encodeChar(char) != "#" {
				builder.WriteRune(char)
				lastIsSpace = false
			}
		}
	}

	return builder.String()
}

func latinToMorse(text string) string {
	builder := strings.Builder{}

	for word := range strings.FieldsSeq(text) {
		if builder.Len() > 0 {
			builder.WriteString(" /")
		}
		for _, char := range word {
			if builder.Len() > 0 {
				builder.WriteRune(' ')
			}
			builder.WriteString(encodeChar(char))
		}
	}

	return builder.String()
}

func morseToLatin(text string) string {
	builder := strings.Builder{}
	lastIsSpace := false

	for word := range strings.SplitSeq(text, "/") {
		if builder.Len() > 0 && !lastIsSpace {
			builder.WriteRune(' ')
			lastIsSpace = true
		}
		for char := range strings.FieldsSeq(word) {
			builder.WriteRune(decodeChar(char))
			lastIsSpace = false
		}
	}

	return builder.String()
}

func morseToMidi(text string, speed byte, preset byte) {
	content := []byte{
		0x4d, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x01,
		0x00, 0x01,
		0x00, speed,
		0x4d, 0x54, 0x72, 0x6b,
		0x00, 0x00, 0x00, 0x1a,
		0x00, 0xb0, 0x07, 0xff,
		0x00, 0xc0, preset,
	}

	for _, char := range text {
		switch char {
		case '.':
			content = append(content,
				0x01, 0x90, 0x45, 0x50,
				0x01, 0x80, 0x45, 0x50,
			)
		case '-':
			content = append(content,
				0x01, 0x90, 0x45, 0x50,
				0x03, 0x80, 0x45, 0x50,
			)
		default:
			content = append(content,
				0x02, 0x80, 0x45, 0x50,
			)
		}
	}

	content = append(content,
		0x02, 0xff, 0x2f, 0x00,
	)

	lenght := fmt.Sprintf("%08x", len(content)-22)
	for idx := range 4 {
		val, _ := strconv.ParseInt(lenght[2*idx:2*idx+2], 16, 0)
		content[idx+18] = byte(val)
	}

	stream, _, _ := midi.Decode(dragonball.NewFile(content), dragonball.NewFont(), dragonball.SampleRate)
	speaker.Play(stream)
}

func Content(minSize fyne.Size) fyne.CanvasObject {
	latinEntry := widget.NewMultiLineEntry()
	latinEntry.SetPlaceHolder("Enter text...")
	latinEntry.Wrapping = fyne.TextWrapWord

	morseEntry := widget.NewMultiLineEntry()
	morseEntry.SetPlaceHolder("Enter text...")
	morseEntry.Wrapping = fyne.TextWrapWord

	entryIsLocked := false

	latinEntry.OnChanged = func(text string) {
		if !entryIsLocked {
			entryIsLocked = true
			morseEntry.SetText(latinToMorse(text))
			entryIsLocked = false
		}
	}

	morseEntry.OnChanged = func(text string) {
		if !entryIsLocked {
			entryIsLocked = true
			latinEntry.SetText(morseToLatin(text))
			entryIsLocked = false
		}
	}

	presetSelect := widget.NewSelect([]string{"Glockenspiel", "Vibraphone"}, nil)
	presetSelect.SetSelected(presetSelect.Options[0])

	speedLabel := widget.NewLabel("Speed:")
	speedSlider := widget.NewSlider(3, 9)
	speedSlider.SetValue(6)

	clearButton := widget.NewButton("Clear", func() {
		latinEntry.SetText(clearLatin(latinEntry.Text))
	})

	playButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		speaker.Clear()
		morseToMidi(morseEntry.Text, byte(speedSlider.Value), presetSelect.Selected[0]-'G'+2)
	})

	stopButton := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
		speaker.Clear()
	})

	return container.New(
		&layout{minSize},
		latinEntry,
		morseEntry,
		clearButton,
		playButton,
		stopButton,
		presetSelect,
		speedLabel,
		speedSlider,
	)
}
