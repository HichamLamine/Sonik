package utils

import (
	// "fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"

	"hichamlamine.github.io/sonik/player"
)

func LoadSongs() []player.Song {
	homedir := os.Getenv("HOME")
	musicdir := filepath.Join(homedir, "Music")

	songsdir, err := os.ReadDir(musicdir)

	if err != nil {
		log.Fatal(err)
	}

	var songsData []player.Song
	for _, file := range songsdir {
		fileExtension := file.Name()[len(file.Name())-3:]
		if !file.IsDir() && fileExtension != "zip" {
			fileDir := filepath.Join(musicdir, file.Name())
			f, err := os.Open(fileDir)
			if err != nil {
				log.Fatal(err)
			}
			m, err := tag.ReadFrom(f)
			if err != nil {
				log.Fatal(err)
			}

			if m.FileType() == "MP3" {
				songsData = append(songsData, player.Song{
					Path:   fileDir,
					Title:  m.Title(),
					Artist: m.Artist(),
					Lyrics: m.Lyrics(),
                    File: f,
				})
			}

		}
	}
	return songsData
}
