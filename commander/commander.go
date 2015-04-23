package commander

import (
	."fmt"
	."../network"
)

func Commander(networkSend chan Message, commanderChan chan Message, aliveChan chan Message, tickerChan chan Message, timerChan chan Message, timeOutChan chan Message, driverOutChan chan Message, driverInChan chan Message) {
	
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
								//Println(message)
								networkSend <- message
							}
						}
						message.Type = "broadcast"
						message.To = 2
						//Println(message)
						networkSend <- message
					}
					notAliveCount++
					
			case <- aliveChan:
				notAliveCount = 0
				//Println("Alive")

			case message = <- commanderChan:
				switch {
				case message.Type == "imAlive" || message.Type == "newElev" || message.Type == "newTarget" || message.Type == "targetUpdate" || message.Type == "addElev" || message.Type == "deleteOrder":
					networkSend <- message

				case message.Type == "signal":
					driverInChan <- message

				case message.Type == "command":	
					if message.Content == "up" {
						message.Type = "stateUpdate"
						message.Content = "MovingUp"
						//Println(message)
						networkSend <- message
						message.Content = "up"
					} else if message.Content == "down" {
						message.Type = "stateUpdate"
						message.Content = "MovingDown"
						//Println(message)
						networkSend <- message
						message.Content = "down"
					} else if message.Content == "stop" {
						message.Type = "door"
						message.Content = "Second"
						message.Value = 3
						timerChan <- message
						message.Type = "stateUpdate"
						message.Content = "Open"
						//Println(message)
						networkSend <- message
						message.Content = "stop"
					}
					message.Type = "engine"
					driverInChan <- message
				}
						
			case message = <- timeOutChan:
				message.Content = "door"
				message.Value = 0
				driverInChan <- message

				message.Type = "stateUpdate"
				message.Content = "Idle"
				//Println(message)
				networkSend <- message
	
			case message = <- driverOutChan:  //floorReached, inside, outsideUp, outsideDown, stop, obstr
				message.Content = message.Type
				message.Floor = message.Floor

				if message.Content == "floorReached" {
					driverInChan <- message
				}

				switch {
				case message.Type == "inside" || message.Type == "outsideUp" || message.Type == "outsideDown":
					message.Content = message.Type
					message.Type = "newOrder"
				case message.Type == "stop" || message.Type == "obstr":
					message.Content = message.Type
					message.Type = "stateUpdate"
				}
				//Println(message)
				networkSend <- message
		}
	}
}