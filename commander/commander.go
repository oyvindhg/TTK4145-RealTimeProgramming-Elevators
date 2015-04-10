package commander

import (
	."fmt"
	."time"
	."../timer"
	."../network"
	."../driver"
)

func Commander(networkSend chan Message, commanderChan chan Message, aliveChan chan Message, tickerChan chan string, timerChan chan TimerInput, timeOutChan chan string, driverInChan chan DriverSignal, driverOutChan chan DriverSignal) {
	notAliveCount := 0
	Sleep(100 * Millisecond)
	timerChan <- TimerInput{150, Millisecond, "alive"}
	
	for {
		select { 							// ADD FLOORREACHED CASE
		case <- tickerChan:
			if notAliveCount == 5 {
				Println("Master dead!")		// IMPLEMENT PANIC
			}
			notAliveCount++
			
		case <- aliveChan:
			notAliveCount = 0
			Println("Alive")

		case command := <- commanderChan:
			switch {
				case command.Content == "imAlive":
					networkSend <- command
				case command.Content == "newID":
					networkSend <- command
				case command.Content == "command":
					if command.Command == "up" {
						driverOutChan <- DriverSignal{"engine", 0, 1}
						Println("going up")
					} else if command.Command == "down" {
						driverOutChan <- DriverSignal{"engine", 0, -1}
					} else if command.Command == "stop" {
						driverOutChan <- DriverSignal{"engine", 0, 0}
					}
				case command.Content == "taskDone":
					Println("taskDone")
				case command.Content == "signal":
					driverOutChan <- DriverSignal{command.ButtonType, command.FloorNumber, 0}
			}
					
		case timeOut := <- timeOutChan:
			Println(timeOut)

		case <- driverInChan:
			//Println(driverIn)
		}
	}
}

// Content = "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", rankChange",
//           "stateUpdate", "connectionChange", "command", "taskDone", "floorReached", "signal"

// RecipientID, SenderID, Content, Command, ElevNumber,
// Online, Rank, FloorNumber, ButtonType, State

// computerID, onlineStatus, rank, floorNum, floorTarget, state, inElev[]

/*
type DriverSignal struct{
	SignalType string  // engine, floorReached, inside, outsideUp, outsideDown, stop, obstr
	FloorNumber int
	Value int
}*/
