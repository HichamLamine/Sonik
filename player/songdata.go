package player

import (
	"log"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type Song struct {
	File   *os.File
	Path   string
	Title  string
	Artist string
	Lyrics string
}

const (
	Paused  bool = false
	Playing bool = true
)

type Player struct {
	playingSong *os.File
	state       bool
	volume      int
	streamer    beep.StreamSeekCloser
}

func (p Player) NewPlayer() Player {
	return Player{}
}

func (p *Player) Play(f *os.File) {
	speaker.Clear()

	p.playingSong = f
	p.state = Playing
	streamer, format, err := mp3.Decode(f)
	speaker.Lock()
	streamer.Seek(0)
	speaker.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	p.streamer = streamer
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(streamer))
}
