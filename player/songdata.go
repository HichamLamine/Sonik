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
	Length time.Duration
	Title  string
	Artist string
	Lyrics string
}

const (
	Paused  bool = false
	Playing bool = true
)

type Player struct {
	PlayingSong *Song
	state       bool
	volume      int
	streamer    beep.StreamSeekCloser
	format      beep.Format
}

func (p Player) NewPlayer() Player {
	return Player{}
}

func (p Player) GetPosition() time.Duration {
	speaker.Lock()
	position := p.format.SampleRate.D(p.streamer.Position())
	speaker.Unlock()
	return position
}

func (p Player) GetLength() time.Duration {
	return p.PlayingSong.Length
}

func (p *Player) Play(s *Song) {
	speaker.Clear()

	p.PlayingSong = s
	p.state = Playing
	streamer, format, err := mp3.Decode(s.File)
	p.PlayingSong.Length = format.SampleRate.D(streamer.Len())
	speaker.Lock()
	_ = streamer.Seek(0)
	speaker.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	p.streamer, p.format = streamer, format
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}
	speaker.Play(beep.Seq(streamer))
}
