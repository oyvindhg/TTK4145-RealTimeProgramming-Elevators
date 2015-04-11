package timer

import (
	."time"
)

type TimerInput struct {
	TimeDuration int
	Scope Duration
	Type string
	ElevNumber int
	RecipientID string
}

func Timekeeper(tickerChan chan string, timerChan chan TimerInput, timeOutChan chan string) {
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

func doorTimer(input TimerInput, timeOutChan chan string) {
	Sleep(Duration(input.TimeDuration) * input.Scope)
	timeOutChan <- input.Type
}

func aliveTicker(input TimerInput, tickerChan chan string) {
	tick := Tick(Duration(input.TimeDuration) * input.Scope)
	for now := range tick {
		tickerChan <- input.Type
		_ = now
	}
}