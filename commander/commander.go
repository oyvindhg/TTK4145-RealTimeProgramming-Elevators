package commander

import (
	."fmt"
	."time"
	."../network"
)

func Commander(networkSend chan Message, commanderChan chan Message, aliveChan chan Message, tickerChan chan string, timerChan chan TimerInput, timeOutChan chan string, driverInChan chan Message, driverOutChan chan Message) {
	
	notAliveCount := 0
	message := Message{}
	Sleep(100 * Millisecond)
	message.Type = "alive"
	message.Content = Millisecond
	message.Value = 150
	timerChan <- message

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

		case commanderMessage := <- commanderChan:
			switch {
			case commanderMessage.Type == "imAlive":
				networkSend <- commanderMessage

			case commanderMessage.Type == "newID":
				networkSend <- commanderMessage

			case commanderMessage.Type == "signal":
				driverOutChan <- commanderMessage

			case commanderMessage.Type == "command":	
				if commanderMessage.Content == "up" {
					message.Value = 1
				} else if commanderMessage.Content == "down" {
					message.Value = -1
				} else if commanderMessage.Content == "stop" {
					message.Type = "door"
					message.Value = 3
					timerChan <- <- message
					message.Value = 0
				}
				message.Type = "engine"
				driverOutChan <- message
			}
					
		case timeOut := <- timeOutChan:
			message.Type = timeOut
			driverOutChan <- message

			message.Type = "stateUpdate"
			message.Content = "Idle"
			networkSend <- message

		case driverInput := <- driverInChan:  //floorReached, inside, outsideUp, outsideDown, stop, obstr
			switch {
			case driverInput.Type == "inside" || driverInput.Type == "outsideUp" || driverInput.Type == "outsideDown":
				message.Type = "newOrder"
			case driverInput.Type == "floorReached" || driverInput.Type == "stop" || driverInput.Type == "obstr":
				message.Type = "stateUpdate"
			}
			message.Content = driverInput.Type
			message.Floor = driverInput.Floor
			networkSend <- message
		}
	}
}