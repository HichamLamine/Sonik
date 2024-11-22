package player

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

type State struct {
	sampleRate    beep.SampleRate
	currentCtrl   *beep.Ctrl
	currentVolume *effects.Volume
}

func NewPlayer() (State, error) {
	sampleRate := beep.SampleRate(48000)
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	return State{sampleRate: sampleRate}, err
}

func decodeAudio(f *os.File) (beep.Streamer, beep.Format, error) {
	switch path.Ext(f.Name()) {
	case ".mp3":
		return mp3.Decode(f)
	case ".wav":
		return wav.Decode(f)
	default:
		return nil, beep.Format{}, fmt.Errorf("unsupported audio format: %s", path.Ext(f.Name()))
	}
}

func (s *State) Play(f *os.File) {
	streamer, format, err := decodeAudio(f)
	if err != nil {
		fmt.Printf("failed to decode file: %v", err)
	}

	// Resample the streamer if the file sample rate isn't the one we use
	if format.SampleRate != s.sampleRate {
		streamer = beep.Resample(4, format.SampleRate, s.sampleRate, streamer)
	}

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	s.currentCtrl = ctrl
	s.currentVolume = volume

    if s.currentVolume != nil {
        speaker.Clear()
    }
	speaker.Play(volume)
}

func (s *State) TogglePause() {
	if s.currentCtrl != nil {
		speaker.Lock()
		s.currentCtrl.Paused = !s.currentCtrl.Paused
		speaker.Unlock()
	}
}

func (s *State) IncreaseVolume() {
	if s.currentVolume != nil {
		speaker.Lock()
		s.currentVolume.Volume += 0.1
		speaker.Unlock()
	}
}

func (s *State) DecreaseVolume() {
	if s.currentVolume != nil {
		speaker.Lock()
		s.currentVolume.Volume -= 0.1
		speaker.Unlock()
	}
}
