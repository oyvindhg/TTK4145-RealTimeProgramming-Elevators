package commander

import (
	."fmt"
	."time"
	."../network"
)

func CommanderInit(networkSend chan Message, commanderChan chan Message, aliveChan chan Message, tickerChan chan Message, timerChan chan Message, timeOutChan chan Message, driverOutChan chan Message, driverInChan chan Message, failureChan chan Message, cancelMasterChan chan Message) {
	
	go commander(commanderChan, networkSend, driverOutChan, driverInChan, timerChan, cancelMasterChan, failureChan)
	go masterAliveHandler(networkSend, tickerChan, timerChan, aliveChan)
	go doorTimeOutHandler(timeOutChan, driverInChan, networkSend)
	go driverOutputHandler(driverOutChan, driverInChan, networkSend, commanderChan)
	go masterChecker(commanderChan)
}
	
func commander(commanderChan chan Message, networkSend chan Message, driverOutChan chan Message, driverInChan chan Message, timerChan chan Message, cancelMasterChan chan Message, failureChan chan Message) {
	for {
		select {
		case message := <- commanderChan:
			//Println("COMMANDER", message)
			switch {
			case message.Type == "findMaster" || message.Type == "newMaster" || message.Type == "addElev" || message.Type == "deleteOrder" || message.Type == "newTarget" || message.Type == "floorReached" || message.Type == "targetUpdate" || message.Type == "floorUpdate":
				networkSend <- message

			case message.Type == "cancelMaster":
				cancelMasterChan <- message

			case message.Type == "masterNumber":
				message.Type = "masterOffline"
				failureChan <- message

			case message.Type == "newOrder":
				driverOutChan <- message

			case message.Type == "signal":
				driverInChan <- message

			case message.Type == "command":	
				if message.Content == "up" {
					message.Type = "stateUpdate"
					message.Content = "MovingUp"
					//Println("\n", message)
					networkSend <- message
					message.Content = "up"
				} else if message.Content == "down" {
					message.Type = "stateUpdate"
					message.Content = "MovingDown"
					//Println("\n", message)
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
				go masterBroadcast(networkSend, cancelMasterChan)
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
			//Println("\n", message)
			networkSend <- message
		}
	}
}

func driverOutputHandler(driverOutChan chan Message, driverInChan chan Message, networkSend chan Message, commanderChan chan Message) {
	for {
		select {
		case message := <- driverOutChan:
			if message.Content == "floorReached" {
				message.Type = message.Content
				driverInChan <- message
			} else if message.Type == "command" {
				commanderChan <- message
				break
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

func masterAliveHandler(networkSend chan Message, tickerChan chan Message, timerChan chan Message, aliveChan chan Message) {
	Sleep(1*Second)
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
				Println("\n", "Master is offline!")	
				networkSend <- message
			}

		case <- aliveChan:
			notAliveCount = 0
		}
	}
}

func masterBroadcast(networkSend chan Message, cancelMasterChan chan Message) {
	message := Message{}
	message.Type = "broadcast"
	message.To = 0
	Println("\n", "Initiating masterBroadcast")
	for {
		select {
		case <- cancelMasterChan:
			Println("\n\n\n\n\nCANCELING MASTER MOAHAHAHAHAHHHH\n\n\n\n\n")
			return
		default:
			networkSend <- message
			Sleep(100 * Millisecond)
		}
	}
}

func masterChecker(commanderChan chan Message) {
	Sleep(1*Second)
	message := Message{}
	message.Type = "newMaster"
	message.Value = 0
	commanderChan <- message
}