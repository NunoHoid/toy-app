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

var letters = [][2]string{
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

var numbers = [][2]string{
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

var markers = [][2]string{
	{"!", "-.-.--"},
	{"\"", ".-..-."},
	{"#", "#"},
	{"$", "#"},
	{"%", "----- -..-. -----"},
	{"&", ".-..."},
	{"'", ".----."},
	{"(", "-.--."},
	{")", "-.--.-"},
	{"*", "-..-"},
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

type applet struct{}

func (_ applet) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(600, 400)
}

func (a applet) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	mediaHeight := float32(40)
	mediaOffset := float32(20)

	entrySize := fyne.NewSize(containerSize.Width, containerSize.Height/2-mediaHeight/2-mediaOffset)

	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(entrySize)

	objects[1].Move(fyne.NewPos(0, entrySize.Height+mediaHeight+2*mediaOffset))
	objects[1].Resize(entrySize)

	start := (containerSize.Width-a.MinSize(nil).Width)/2 + mediaOffset
	end := start + a.MinSize(nil).Width - 2*mediaOffset

	objects[2].Move(fyne.NewPos(start, entrySize.Height+mediaOffset))
	objects[2].Resize(fyne.NewSquareSize(mediaHeight))

	objects[3].Move(objects[2].Position().AddXY(mediaHeight+mediaOffset, 0))
	objects[3].Resize(fyne.NewSquareSize(mediaHeight))

	objects[4].Move(objects[3].Position().AddXY(mediaHeight+mediaOffset, 0))
	objects[4].Resize(fyne.NewSize(objects[4].MinSize().Width, mediaHeight))

	objects[5].Move(objects[4].Position().AddXY(objects[4].Size().Width+mediaOffset, 0))
	objects[5].Resize(fyne.NewSize(objects[5].MinSize().Width, mediaHeight))

	objects[6].Move(objects[5].Position().AddXY(objects[5].Size().Width, 0))
	objects[6].Resize(fyne.NewSize(end-objects[6].Position().X, mediaHeight))
}

func encodeChar(char rune) string {
	if 'A' <= char && char <= 'Z' {
		return letters[char-'A'][1]
	}
	if 'a' <= char && char <= 'z' {
		return letters[char-'a'][1]
	}
	if '0' <= char && char <= '9' {
		return numbers[char-'0'][1]
	}
	if '!' <= char && char <= '/' {
		return markers[char-'!'][1]
	}
	if ':' <= char && char <= '@' {
		return markers[char-':'][1]
	}
	return "#"
}

func decodeChar(char string) rune {
	if idx := slices.IndexFunc(letters, func(pair [2]string) bool { return char == pair[1] }); idx != -1 {
		return rune(letters[idx][0][0])
	}
	if idx := slices.IndexFunc(numbers, func(pair [2]string) bool { return char == pair[1] }); idx != -1 {
		return rune(numbers[idx][0][0])
	}
	if idx := slices.IndexFunc(markers, func(pair [2]string) bool { return char == pair[1] }); idx != -1 {
		return rune(markers[idx][0][0])
	}
	return '#'
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
		0x00, 0x00,
		0x00, 0x01,
		0x00, speed,
		0x4d, 0x54, 0x72, 0x6b,
		0x00, 0x00, 0x00, 0x04,
		0x00, 0xb0, 0x07, 0xff,
		0x00, 0xc0, preset,
	}

	for _, char := range text {
		switch char {
		case '.':
			content = append(content,
				0x01, 0x90, 0x45, 0x40,
				0x01, 0x80, 0x45, 0x40,
			)
		case '-':
			content = append(content,
				0x01, 0x90, 0x45, 0x40,
				0x03, 0x80, 0x45, 0x40,
			)
		default:
			content = append(content,
				0x02, 0x80, 0x45, 0x40,
			)
		}
	}

	content = append(content,
		0x07, 0x80, 0x45, 0x40,
		0x00, 0xff, 0x2f, 0x00,
	)

	lenght := fmt.Sprintf("%08x", len(content)-22)
	for idx := range 4 {
		val, _ := strconv.ParseInt(lenght[2*idx:2*idx+2], 16, 0)
		content[idx+18] = byte(val)
	}

	stream, _, _ := midi.Decode(dragonball.NewFile(content), dragonball.NewFont(), dragonball.SampleRate)
	speaker.Play(stream)
}

func GuiContent() fyne.CanvasObject {
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

	speedLabel := widget.NewLabel("Playback speed:")

	speedSlider := widget.NewSlider(1, 9)
	speedSlider.SetValue(5)

	presetSelect := widget.NewSelect([]string{"Glockenspiel", "Vibraphone"}, nil)
	presetSelect.SetSelected(presetSelect.Options[0])

	playButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		speaker.Clear()
		morseToMidi(morseEntry.Text, byte(speedSlider.Value), presetSelect.Selected[0]-'G'+2)
	})

	stopButton := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
		speaker.Clear()
	})

	return container.New(
		&applet{},
		latinEntry,
		morseEntry,
		playButton,
		stopButton,
		presetSelect,
		speedLabel,
		speedSlider,
	)
}
