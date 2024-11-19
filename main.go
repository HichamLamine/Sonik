package main

import (
	"os"

	"github.com/rivo/tview"
	"hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/utils"
	// "hichamlamine.github.io/sonik/model"
)

func main() {
	app := tview.NewApplication()
	list := tview.NewList()

	player := player.Player{}

	songs := utils.LoadSongs()

	for _, song := range songs {
		list.AddItem(song.Title, song.Artist, ' ', nil)
	}
	list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		var songFile *os.File
		for _, song := range songs {
			if song.Title == s1 {
				songFile = song.File
			}
		}
		go player.Play(songFile)
	})

	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}

}
