package dragonball

import (
	_ "embed"
	"io"
	"time"

	"github.com/gopxl/beep/v2/midi"
	"github.com/gopxl/beep/v2/speaker"
)

type file struct {
	content []byte
}

func (f *file) Read(slice []byte) (int, error) {
	if len(f.content) == 0 {
		return 0, io.EOF
	}
	n := copy(slice, f.content)
	f.content = f.content[n:]
	return n, nil
}

func (f *file) Close() error {
	return nil
}

//go:embed dragonball.sf2
var dragonball []byte

const SampleRate = 48000

func NewFile(content []byte) *file {
	return &file{content}
}

func NewFont() *midi.SoundFont {
	font, _ := midi.NewSoundFont(NewFile(dragonball))
	return font
}

func SpeakerInit() {
	content := []byte{
		0x4d, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x00,
		0x00, 0x01,
		0x00, 0x05,
		0x4d, 0x54, 0x72, 0x6b,
		0x00, 0x00, 0x00, 0x04,
		0x00, 0xff, 0x2f, 0x00,
	}

	_, format, _ := midi.Decode(NewFile(content), NewFont(), SampleRate)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}
