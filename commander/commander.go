package commander

import (
	."fmt"
	."../network"
)

func Commander(networkSend chan Message, commanderChan chan Message, aliveChan chan Message, tickerChan chan string, timerChan chan Message, timeOutChan chan string, driverInChan chan Message, driverOutChan chan Message) {
	
	notAliveCount := 0
	message := Message{}
	message.Type = "alive"
	message.Content = "Millisecond"
	message.Value = 150
	timerChan <- message

	for {
		select {
			case <- tickerChan:
					if notAliveCount == 5 {
						Println("Master dead!")		// IMPLEMENT PANIC
						message.Type = "elevOffline"
						message.From = 1
						for i := 2; i < ELEV_COUNT + 1; i++ {
							if i != destination {		
								Println("Sending elevOffline message to elev", i)				
								message.To = i
								networkSend <- message
							}
						}
						message = "broadcast"
						message.To = 2
						networkSend <- message
					}
					notAliveCount++
					
				case <- aliveChan:
					notAliveCount = 0
					Println("Alive")
	
				case commanderMessage := <- commanderChan:
					switch {
					case commanderMessage.Type == "imAlive" || commanderMessage.Type == "newElev" || commanderMessage.Type == "addElev":
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
							timerChan <- message
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