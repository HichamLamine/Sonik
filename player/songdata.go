package player

import (
	"os"
	"time"
)

type Song struct {
	File   *os.File
	Path   string
	Length time.Duration
	Title  string
	Artist string
	Lyrics string
}
