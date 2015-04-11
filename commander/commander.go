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
	masterID := ""+ PORT
	elevRank := 0
	Sleep(100 * Millisecond)
	timerChan <- TimerInput{150, Millisecond, "alive", 0, ""}
	
	for {
		select {
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
				} else if command.Command == "down" {
					driverOutChan <- DriverSignal{"engine", 0, -1}
				} else if command.Command == "stop" {
					driverOutChan <- DriverSignal{"engine", 0, 0}
					timerChan <- TimerInput{3, Second, "door", command.ElevNumber, command.RecipientID}
				}

			case command.Content == "taskDone":
				Println("taskDone")
				
			case command.Content == "signal":
				driverOutChan <- DriverSignal{command.ButtonType, command.FloorNumber, command.Rank}
			}
					
		case timeOut := <- timeOutChan:
			networkSend <- Message{timeOut.RecipientID, "", "stateUpdate", timeOut.ElevNumber, true, 0, 0, "", "Idle"}
			driverOutChan <- DriverSignal{"door", 0, 1}

		case driverInput := <- driverInChan:  //floorReached, inside, outsideUp, outsideDown, stop, obstr
			switch {
			case driverInput.SignalType == "inside" || driverInput.SignalType == "outsideUp" || driverInput.SignalType == "outsideDown":
				networkSend <- Message{RECIPIENTID, "", "newOrder", "", 0, true, 0, driverInput.FloorNumber, driverInput.SignalType, ""}
			case driverInput.SignalType == "floorReached" || driverInput.SignalType == "stop" || driverInput.SignalType == "obstr":
				networkSend <- Message{RECIPIENTID, "", "stateUpdate", "", 0, true, 0, driverInput.FloorNumber, driverInput.SignalType, ""}
			}
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
