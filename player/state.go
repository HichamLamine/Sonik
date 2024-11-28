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
	sampleRate      beep.SampleRate
	currentStreamer beep.StreamSeekCloser
	currentCtrl     *beep.Ctrl
	currentVolume   *effects.Volume
	PlayingSong     *Song
	universalVolume float64
}

func NewPlayer() (State, error) {
	sampleRate := beep.SampleRate(48000)
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	return State{sampleRate: sampleRate, universalVolume: 0.5}, err
}

func decodeAudio(f *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	switch path.Ext(f.Name()) {
	case ".mp3":
		return mp3.Decode(f)
	case ".wav":
		return wav.Decode(f)
	default:
		return nil, beep.Format{}, fmt.Errorf("unsupported audio format: %s", path.Ext(f.Name()))
	}
}

func (s *State) Play(song *Song) {
	// reset the file pointer so it wouldn't try to read from the end of it
	_, err := song.File.Seek(0, 0)
	if err != nil {
		fmt.Println("failed to reset file pointer to the beginning")
		os.Exit(1)
	}

	streamer, _, err := decodeAudio(song.File)
	if err != nil {
		fmt.Printf("failed to decode file: %v", err)
	}

	loop := beep.Loop(-1, streamer)

	// Resample the streamer if the file sample rate isn't the one we use
	// if format.SampleRate != s.sampleRate {
	// 	streamer = beep.Resample(4, format.SampleRate, s.sampleRate, streamer)
	// }

	ctrl := &beep.Ctrl{Streamer: loop, Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	if s.currentVolume != nil {
		speaker.Clear()
	}

	s.PlayingSong = song
	s.currentStreamer = streamer
	s.currentCtrl = ctrl
	s.currentVolume = volume

	s.SetVolume(Denormalize(s.universalVolume))
	if s.universalVolume == 0 {
		s.ToggleMute(true)
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
		s.currentVolume.Silent = false
		if s.Normalized() >= 0.99 {
			s.currentVolume.Volume = Denormalize(1)
		} else {
			s.currentVolume.Volume = Denormalize(s.Normalized() + 0.1)
		}
		speaker.Unlock()
	}
	s.universalVolume = s.Normalized()
}

func (s *State) DecreaseVolume() {
	if s.currentVolume != nil {
		speaker.Lock()
		if s.Normalized() < 0.09 {
			s.currentVolume.Volume = Denormalize(0.)
		} else {
			s.currentVolume.Volume = Denormalize(s.Normalized() - 0.1)
		}
		if s.Normalized() == 0 {
			s.currentVolume.Silent = true
		}
		speaker.Unlock()
	}
	s.universalVolume = s.Normalized()
}

func (s *State) SetVolume(v float64) {
	if s.currentVolume != nil {
		speaker.Lock()
		s.currentVolume.Volume = v
		speaker.Unlock()
	}
}

func (s *State) ToggleMute(b bool) {
	if s.currentVolume != nil {
		speaker.Lock()
		s.currentVolume.Silent = b
		speaker.Unlock()
	}
}

func (s State) GetPosition() time.Duration {
	if s.currentStreamer == nil {
		return time.Second * 0
	} else {
		return s.sampleRate.D(s.currentStreamer.Position())
	}
}

func (s State) GetPositionPercent() float64 {
	return float64(s.GetPosition()) / float64(s.GetLen())
}

func (s State) GetLen() time.Duration {
	if s.currentStreamer == nil {
		return time.Second * 0
	} else {
		return s.sampleRate.D(s.currentStreamer.Len())
	}
}

func (s State) IsPaused() bool {
	if s.currentCtrl == nil {
		return true
	}
	return s.currentCtrl.Paused
}

func (s State) Normalized() float64 {
	min := -6.7
	midPoint := -3.
	max := -2.7

	v := s.currentVolume.Volume
	if v <= -3 {
		return 0.5 * (v - min) / (midPoint - min)
	} else {
		return 0.5 + 0.5*(v-midPoint)/(max-midPoint)
	}
}

// (-4.8 -> 1.8) => 0.0 -> 1.0 (-4.8, -2.1) -> (0, 0.5), (-2.1 -> 1.8) -> (0.5, 1.0)

func Denormalize(v float64) float64 {
	min := -6.7
	midPoint := -3.
	max := -2.7

	if v <= 0.5 {
		return min + (v/0.5)*(midPoint-min)
	} else {
		return midPoint + (v-0.5)/0.5*(max-midPoint)
	}
}

func (s State) GetVolume() float64 {
	return s.universalVolume
}

func (s State) SeekForward() {
	if s.currentStreamer != nil {
		speaker.Lock()
		s.currentCtrl.Paused = true
		position := s.GetPosition()
		newPosition := s.sampleRate.N(position + 5*time.Second)
		if newPosition <= s.currentStreamer.Len() {
			_ = s.currentStreamer.Seek(newPosition)
		} else {
			_ = s.currentStreamer.Seek(s.currentStreamer.Len())
		}
		s.currentCtrl.Paused = false
		speaker.Unlock()
	}
}

func (s State) SeekBackward() {
	if s.currentCtrl.Streamer != nil {
		speaker.Lock()
		position := s.GetPosition()
		newPosition := s.sampleRate.N(position - 5*time.Second)
		if newPosition >= 0 {
			_ = s.currentStreamer.Seek(newPosition)
		} else {
			_ = s.currentStreamer.Seek(0)
		}
		speaker.Unlock()
	}
}
