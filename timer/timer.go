package timer

import (
	//."fmt"
	."time"
	."../network"
)

func Timekeeper(tickerChan chan Message, timerChan chan Message, timeOutChan chan Message) {
	closeDoor := 0
	for {
		select {
		case message := <- timerChan:
			switch {
			case message.Type == "door":
				closeDoor++
				go doorTimer(&closeDoor, message, timeOutChan)
			case message.Type == "alive":
				go aliveTicker(message, tickerChan)
			}
		}
	}
}

func doorTimer(closeDoor *int, message Message, timeOutChan chan Message) {
	switch{
	case message.Content == "Second":
		Sleep(Duration(message.Value) * Second)
	case message.Content == "Millisecond":
		Sleep(Duration(message.Value) * Millisecond)
	case message.Content == "MicroSecond":
		Sleep(Duration(message.Value) * Microsecond)
	}
	*closeDoor--
	if *closeDoor == 0 {
		timeOutChan <- message
	}
}

func aliveTicker(message Message, tickerChan chan Message) {
	tick := Tick(0 * Second)
	switch{
		case message.Content == "Second":
			tick = Tick(Duration(message.Value) * Second)
		case message.Content == "Millisecond":
			tick = Tick(Duration(message.Value) * Millisecond)
		case message.Content == "Microsecond":
			tick = Tick(Duration(message.Value) * Microsecond)
	}
	for now := range tick {
		tickerChan <- message
		_ = now
	}
}