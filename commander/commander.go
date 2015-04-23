package commander

import (
	."fmt"
	."../network"
)

func Commander(networkSend chan Message, commanderChan chan Message, aliveChan chan Message, tickerChan chan Message, timerChan chan Message, timeOutChan chan Message, driverInChan chan Message, driverOutChan chan Message) {
	
	notAliveCount := 0
	message := Message{}
	message.Type = "alive"
	message.Content = "Millisecond"
	message.Value = 150
	timerChan <- message

	for {
		select {
			case message = <- tickerChan:
					if notAliveCount == 5 {
						destination := message.To
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
						message.Type = "broadcast"
						message.To = 2
						networkSend <- message
					}
					notAliveCount++
					
			case <- aliveChan:
				notAliveCount = 0
				//Println("Alive")

			case commanderMessage := <- commanderChan:
				switch {
				case commanderMessage.Type == "imAlive" || commanderMessage.Type == "newElev" || commanderMessage.Type == "newTarget" || commanderMessage.Type == "targetUpdate" || commanderMessage.Type == "addElev" || commanderMessage.Type == "deleteOrder":
					networkSend <- commanderMessage

				case commanderMessage.Type == "signal":
					driverOutChan <- commanderMessage

				case commanderMessage.Type == "command":	
					if commanderMessage.Content == "up" {
						message.Type = "stateUpdate"
						message.Content = "MovingUp"
						networkSend <- message
						message.Content = "up"
					} else if commanderMessage.Content == "down" {
						message.Type = "stateUpdate"
						message.Content = "MovingDown"
						networkSend <- message
						message.Content = "down"
					} else if commanderMessage.Content == "stop" {
						message.Type = "door"
						message.Content = "Second"
						message.Value = 3
						timerChan <- message
						message.Type = "stateUpdate"
						message.Content = "Open"
						networkSend <- message
						message.Content = "stop"
					}
					message.Type = "engine"
					driverOutChan <- message
				}
						
			case message = <- timeOutChan:
				message.Content = "door"
				message.Value = 0
				driverOutChan <- message

				message.Type = "stateUpdate"
				message.Content = "Idle"
				networkSend <- message
	
			case driverInput := <- driverInChan:  //floorReached, inside, outsideUp, outsideDown, stop, obstr
				driverInput.Content = driverInput.Type
				driverInput.Floor = driverInput.Floor
				switch {
				case driverInput.Type == "inside" || driverInput.Type == "outsideUp" || driverInput.Type == "outsideDown":
					driverInput.Content = driverInput.Type
					driverInput.Type = "newOrder"
				case driverInput.Type == "stop" || driverInput.Type == "obstr":
					driverInput.Content = driverInput.Type
					driverInput.Type = "stateUpdate"
				}
				networkSend <- driverInput
		}
	}
}