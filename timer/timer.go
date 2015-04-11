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
	switch{
	case input.Content == "Second":
		Sleep(Duration(input.Value) * Second)
	case input.Content == "Millisecond":
		Sleep(Duration(input.Value) * Millisecond)
	case input.Content == "MicroSecond":
		Sleep(Duration(input.Value) * Microsecond)
	}
	timeOutChan <- input.Type
}

func aliveTicker(input Message, tickerChan chan string) {
	tick := Tick(0 * Second)
	switch{
		case input.Content == "Second":
			tick = Tick(Duration(input.Value) * Second)
		case input.Content == "Millisecond":
			tick = Tick(Duration(input.Value) * Millisecond)
		case input.Content == "Microsecond":
			tick = Tick(Duration(input.Value) * Microsecond)
	}
	for now := range tick {
		tickerChan <- input.Type
		_ = now
	}
}