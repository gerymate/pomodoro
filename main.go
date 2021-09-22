package main

/* A simple cross-platform pomodoro timer sitting in the systray */

import (
	"gerymate/pomodoro/icon"
	"fmt"
	"time"
	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Pomodoro Timer")
	systray.SetTooltip("Limit yourself to stay efficient")
	mPom := systray.AddMenuItem("Start pomodoro", "Start a 25 minutes focus period")
	mBreak := systray.AddMenuItem("Start break", "Start a short break")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Goodbye!")
		systray.Quit()
	}()

	go func() {
		for {
			select {
				case <-mPom.ClickedCh:
					fmt.Println("Clicked Pomodoro")
					pomodoro("Pomodoro", "Focus...", 25*60)
				case <-mBreak.ClickedCh:	
					fmt.Println("Clicked Break")
					pomodoro("Pomodoro", "Take a break!", 5*60)
			}
		}
	}()
}

func pomodoro(title string, text string, length int) {
	dlg, err := zenity.Progress(
		zenity.Title(title))
	if err != nil {
		return
	}
	defer dlg.Close()
	
	dlg.Text(text)
	
	is_closed := dlg.Done()

	for elapsed := 0; elapsed < length; elapsed++ {
		select {
		case _, ok := <-is_closed:
				if !ok {
					return
				}
			default:
				var value float32 = float32(elapsed) / float32(length) * 100
				dlg.Value(int(value))
				time.Sleep(time.Second)
		}
	}

	dlg.Complete()
	time.Sleep(time.Second)
}

func onExit() {
	// clean up here
}