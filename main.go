package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"hichamlamine.github.io/sonik/player"
	"hichamlamine.github.io/sonik/utils"
	// "hichamlamine.github.io/sonik/model"
)

const (
	playIcon     string = ""
	pauseIcon    string = "󰏤"
	nextIcon     string = ""
	previousIcon string = ""
)

func main() {
	app := tview.NewApplication()

	player := player.Player{}

	songs := utils.LoadSongs()

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Songs")
	for _, song := range songs {
		list.AddItem(song.Title, song.Artist, ' ', nil)
	}
	list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		for _, song := range songs {
			if song.Title == s1 {
				go player.Play(&song)
			}
		}
	})

	var position time.Duration
	var length time.Duration
	progressView := tview.NewTextView().SetChangedFunc(func() { app.Draw() })
	progressView.SetBorder(true)
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				if player.PlayingSong == nil {
					position, length = 0, 0
				} else {
					progressView.Clear()
					position, length = player.GetPosition(), player.GetLength()
					progressText := fmt.Sprintf("\n%v/%v\n%s  %s  %s", position.Round(time.Second), length.Round(time.Second), previousIcon, playIcon, nextIcon)
					progressView.Write([]byte(progressText))
				}
			}
		}
	}()

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true).
		AddItem(progressView.SetTextAlign(tview.AlignCenter), 5, 1, false)

	if err := app.SetRoot(flex, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}

}
