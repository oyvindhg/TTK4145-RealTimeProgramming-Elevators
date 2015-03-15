package timer

import (
	."time"
)

type TimerInput struct {
	TimeDuration int
	Scope Duration
	Type string
}

func InitTimer(tickerChan chan string, timerChan chan TimerInput, timeOutChan chan string) {
	for {
		select {
		case input := <- timerChan:
			if input.Type == "door" {
				go timer(input, timeOutChan)
			} else if input.Type == "alive"{
				go ticker(input, tickerChan)
			}
		}
	}
}

func timer(input TimerInput, timeOutChan chan string) {

	Sleep(Duration(input.TimeDuration) * input.Scope)
	
	timeOutChan <- input.Type
}

func ticker(input TimerInput, tickerChan chan string) {
	tick := Tick(Duration(input.TimeDuration) * input.Scope)
	for now := range tick {
		tickerChan <- input.Type
		_ = now
	}
	
}