package liftState

import (
	."fmt"
	."time"
	."../network"
)

type elevator struct {
	floorNum int
	floorTarget int
	state string
}

func LiftStateInit(networkReceive chan Message, commanderChan chan Message, aliveChan chan Message, fileOutChan chan Message, fileInChan chan Message) {

	message := Message{}
	elev := make([]elevator, 1, ELEV_COUNT + 1)
	inside 	:= make([]int, FLOOR_COUNT+1)
	outUp 	:= make([]int, FLOOR_COUNT+1)
	outDown	:= make([]int, FLOOR_COUNT+1)

	message.Type = "newElev"
	commanderChan <- message
	for i := 1; i < FLOOR_COUNT + 1; i++ {
		message.Type = "readInside"
		message.Floor = i
		fileInChan <- message
		message = <- fileOutChan
		if message.Value != -1 {
			inside[i] = message.Value
			if message.Value == 1 {
				message.Type = "newOrder"

				message.Content = "inside"
				commanderChan <- message
			}
		}
	}

	go liftState(networkReceive, commanderChan, aliveChan)
}

func liftState(networkReceive chan Message, commanderChan chan Message, aliveChan chan Message) {
	for {
		select {
		case message = <- networkReceive:
			switch {
			case message.Type == "command" || message.Type == "master":
				Println(message)
				commanderChan <- message

			case message.Type == "imAlive":
				aliveChan <- message

			case message.Type == "newElev" || message.Type == "addElev":
				Println(message)
				elev = append(elev, elevator{0, 0, "Idle"})

			case message.Type == "elevOffline":
				Println(message)
				elev = append(elev[:message.Value], elev[message.Value+1:]...)

			case message.Type == "newOrder":
				Println(message)
				if message.To == 1 {
					bestFloor := FLOOR_COUNT
					bestElev := 0
					for i := 1; i < len(elev); i++ {
						if message.Floor == FLOOR_COUNT && elev[i].floorNum == FLOOR_COUNT-1 && elev[i].state == "MovingUp"{
							break
						} else if message.Floor == 1 && elev[i].floorNum == 2 && elev[i].state == "MovingDown"{
							break
						} else if message.Floor - elev[i].floorNum == 1 && elev[i].state == "MovingUp"{
							break
						} else if message.Floor - elev[i].floorNum == -1 && elev[i].state == "MovingDown"{
							break
						}
						if elev[i].state == "Idle" || elev[i].state == "Open" {
							if message.Floor - elev[i].floorNum >= 0 && message.Floor - elev[i].floorNum < bestFloor {
								bestFloor = message.Floor - elev[i].floorNum
								bestElev = i
							} else if elev[i].floorNum - message.Floor > 0 && elev[i].floorNum - message.Floor < bestFloor {
								bestFloor = elev[i].floorNum - message.Floor
								bestElev = i
							}
						}	
					}
					message.To = bestElev
					message.Type = "newTarget"
					commanderChan <- message
				}

				switch{
				case message.Content == "inside":
					inside[message.Floor] = 1
					message.Type = "writeInside"
					message.Value = 1
					fileInChan <- message
				case message.Content == "outsideUp":
					outUp[message.Floor] = 1
				case message.Content == "outsideDown":
					outDown[message.Floor] = 1
				}
				message.Type = "signal"
				message.Value = 1
				commanderChan <- message


			case message.Type == "deleteOrder":
				Println(message)
				switch{
				case message.Content == "inside":
					inside[message.Floor] = 0
					message.Type = "writeInside"
					message.Value = 0
					fileInChan <- message
				case message.Content == "outsideUp":
					outUp[message.Floor] = 0
				case message.Content == "outsideDown":
					outDown[message.Floor] = 0
				}
				
				if message.Content == "inside" && message.From != message.To {
					break
				}
				message.Type = "signal"
				message.Value = 0
				commanderChan <- message

			case message.Type == "floorUpdate":
				Println(message)
				elev[message.From].floorNum = message.Floor

			case message.Type == "floorReached":
				Println(message)
				elev[message.From].floorNum = message.Floor
				if elev[message.From].floorTarget == message.Floor || inside[message.Floor] == 1 || outUp[message.Floor] == 1 && elev[message.From].state == "MovingUp" || outDown[message.Floor] == 1 && elev[message.From].state == "MovingDown"{
					message.Type = "command"
					message.Content = "stop"
					commanderChan <- message
					message.Type = "deleteOrder"
					if inside[message.Floor] == 1 {
						message.Content = "inside"
						commanderChan <- message
					}
					if outUp[message.Floor] == 1 {
						message.Content = "outsideUp"
						commanderChan <- message
					}
					if outDown[message.Floor] == 1 {
						message.Content = "outsideDown"
						commanderChan <- message
					}
					if elev[message.From].floorTarget == message.Floor {	
						message.Type = "targetUpdate"
						message.Floor = 0
						commanderChan <- message
					}
					break
				}
				emptyQueue := true
				for i := 1; i < FLOOR_COUNT + 1; i++ {
					if inside[i] == 1 || outDown[i] == 1 || outUp[i] == 1 {
						emptyQueue = false
					}
				}
				if emptyQueue == true {
					message.Type = "command"
					message.Content = "stop"
					commanderChan <- message
				}

			case message.Type == "newTarget":
				Println(message)
				message.Type = "targetUpdate"
				commanderChan <- message
				message.Type = "command"
				if message.Floor > elev[message.To].floorNum && elev[message.To].state != "Open" {
					message.Content = "up"
					commanderChan <- message
				} else if message.Floor < elev[message.To].floorNum && elev[message.To].state != "Open" {
					message.Content = "down"
					commanderChan <- message
				} else if message.Floor == elev[message.To].floorNum {
					message.Type = "deleteOrder"
					commanderChan <- message
					message.Type = "command"
					message.Content = "stop"
					commanderChan <- message
				}

			case message.Type == "targetUpdate":
				Println(message)
				elev[message.From].floorTarget = message.Floor

			case message.Type == "stateUpdate":
				elev[message.From].state = message.Content
				insideQueueEmpty := true
				if elev[message.From].state == "Idle" && message.To == message.From {
					for i := 1; i < FLOOR_COUNT + 1; i++ {
						if elev[message.From].floorNum + i < FLOOR_COUNT + 1 {
						 	if inside[elev[message.From].floorNum + i] == 1  {
						 		insideQueueEmpty = false
								message.Type = "newTarget"
								message.Floor = elev[message.From].floorNum + i
								commanderChan <- message
								break
							}
						}
						if elev[message.From].floorNum - i > 0 {
							if inside[elev[message.From].floorNum - i] == 1 {
								insideQueueEmpty = false
								message.Type = "newTarget"
								message.Floor = elev[message.From].floorNum - i
								commanderChan <- message
								break
							}
						}
					}
				}
				if insideQueueEmpty == false {
					break
				}

				if elev[message.From].state == "Idle" && message.To == 1 {
					for i := 1; i < FLOOR_COUNT + 1; i++ {
						if elev[message.From].floorNum + i < FLOOR_COUNT + 1 {
						 	if outDown[elev[message.From].floorNum + i] == 1 ||  outUp[elev[message.From].floorNum + i] == 1 {
								message.To = message.From
								message.Type = "newTarget"
								message.Floor = elev[message.From].floorNum + i
								commanderChan <- message
								break
							}
						}
						if elev[message.From].floorNum - i > 0 {
							if outDown[elev[message.From].floorNum - i] == 1 ||  outUp[elev[message.From].floorNum - i] == 1 {
								message.To = message.From
								message.Type = "newTarget"
								message.Floor = elev[message.From].floorNum - i
								commanderChan <- message
								break
							}
						}
					}
				}
			}
		}
	}
}