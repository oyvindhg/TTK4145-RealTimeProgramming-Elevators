package commander

import (
	."fmt"
	."time"
	."../network"
)

func CommanderInit(networkSend chan Message, commanderChan chan Message, aliveChan chan Message, tickerChan chan Message, timerChan chan Message, timeOutChan chan Message, driverOutChan chan Message, driverInChan chan Message, failureChan chan Message) {
	
	go commander(commanderChan, networkSend, driverOutChan, driverInChan, timerChan)
	go masterAliveHandler(tickerChan, timerChan, aliveChan, failureChan)
	go doorTimeOutHandler(timeOutChan, driverInChan, networkSend)
	go driverOutputHandler(driverOutChan, driverInChan, networkSend)
}
	
func commander(commanderChan chan Message, networkSend chan Message, driverOutChan chan Message, driverInChan chan Message, timerChan chan Message) {
	for {
		select {
		case message := <- commanderChan:
			switch {
			case message.Type == "findMaster" || message.Type == "newTarget" || message.Type == "floorReached" || message.Type == "targetUpdate" || message.Type == "floorUpdate" || message.Type == "addElev" || message.Type == "deleteOrder":
				networkSend <- message

			case message.Type == "newOrder":
				driverOutChan <- message

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
					networkSend <- message
					message.Content = "stop"
				}
				message.Type = "engine"
				driverInChan <- message

			case message.Type == "master":
				go masterAliveBroadcast(networkSend)
			}
		}
	}
}

func doorTimeOutHandler(timeOutChan chan Message, driverInChan chan Message, networkSend chan Message) {
	for {
		select {
		case message := <- timeOutChan:
			message.Content = "door"
			message.Value = 0
			driverInChan <- message

			message.Type = "stateUpdate"
			message.Content = "Idle"
			//Println(message)
			networkSend <- message
		}
	}
}

func driverOutputHandler(driverOutChan chan Message, driverInChan chan Message, networkSend chan Message) {
	for {
		select {
		case message := <- driverOutChan:
			if message.Content == "floorReached" {
				message.Type = message.Content
				driverInChan <- message
			}
			switch {
			case message.Content == "inside" || message.Content == "outsideUp" || message.Content == "outsideDown":
				message.Type = "newOrder"
			case message.Content == "stopButton" || message.Content == "obstrOn" || message.Content == "obstrOff":
				message.Type = "stateUpdate"
			}
			networkSend <- message
		}
	}
}

func masterAliveHandler(tickerChan chan Message, timerChan chan Message, aliveChan chan Message, failureChan chan Message) {
	notAliveCount := 0
	message := Message{}
	message.Type = "alive"
	message.Content = "Millisecond"
	message.Value = 150
	timerChan <- message
	for {
		select {
		case message = <- tickerChan:
			notAliveCount++
			if notAliveCount == 5 {
				message.Type = "masterOffline"
				Println("Master is offline!")	
				failureChan <- message
			}

		case <- aliveChan:
			notAliveCount = 0
		}
	}
}

func masterAliveBroadcast(networkSend chan Message) {
	message := Message{}
	message.Type = "lookForElevators"
	networkSend <- message
	message.Type = "imAlive"
	message.To = 0
	Println("Initiating masterAliveBroadcast")
	for {
		Sleep(100 * Millisecond)
		networkSend <- message
	}
}
