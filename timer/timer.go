package timer

import (
	."time"
	."../network"
)

func Timekeeper(tickerChan chan string, timerChan chan Message, timeOutChan chan string) {
	for {
		select {
		case input := <- timerChan:
			switch {
			case input.Type == "door":
				go doorTimer(input, timeOutChan)
			case input.Type == "alive":
				go aliveTicker(input, tickerChan)
			}
		}
	}
}

func doorTimer(input Message, timeOutChan chan string) {
	Sleep(Duration(input.Value) * input.Content)
	timeOutChan <- input.Type
}

func aliveTicker(input Message, tickerChan chan string) {
	tick := Tick(Duration(input.Value) * input.Content)
	for now := range tick {
		tickerChan <- input.Type
		_ = now
	}
}