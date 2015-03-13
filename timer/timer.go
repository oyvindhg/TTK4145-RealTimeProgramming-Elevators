package timer

import (
	."time"
)

type TimerInput struct {
	TimeDuration int
	Scope string
	Type string
}

func InitTimer(timerChan chan TimerInput, timeOutChan chan string) {
	for {
		select {
		case input := <- timerChan:
			if input.Type == "doorOpen" {
				go timer(input, timeOutChan)
			} else {
				go ticker(input, timeOutChan)
			}
		}
	}
}

func timer(input TimerInput, timeOutChan chan string) {
	switch {
	case input.Scope == "Second":
		Sleep(Duration(input.TimeDuration) * Second)
	case input.Scope == "Millisecond":
		Sleep(Duration(input.TimeDuration) * Millisecond)
	case input.Scope == "Microsecond":
		Sleep(Duration(input.TimeDuration) * Microsecond)
	}
	timeOutChan <- input.Type
}

func ticker(input TimerInput, timeOutChan chan string) {
	
}