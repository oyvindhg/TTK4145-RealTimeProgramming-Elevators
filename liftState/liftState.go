package liftState

import (
	//."fmt"
	."time"
	."../network"
	//."../fileManager"
)

type elevator struct {
	floorNum int
	floorTarget int
	state string
}

func aliveBroadcast(commanderChan chan Message) {
	message := Message{}
	message.Type = "imAlive"
	for {	
		Sleep(100 * Millisecond)
		commanderChan <- message
	}
}

func LiftState(networkReceive chan Message, commanderChan chan Message, aliveChan chan Message, fileInChan chan Message, fileOutChan chan Message) {

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
		fileOutChan <- message
		message = <- fileInChan
		if message.Value != -1 {
			inside[i] = message.Value
		}
	}
	
	for{
		select{
		case message = <- networkReceive:
			switch{
			case message.Type == "master":
				go aliveBroadcast(commanderChan)

			case message.Type == "imAlive":
				aliveChan <- message

			case message.Type == "command":
				commanderChan <- message

			case message.Type == "newElev" || message.Type == "addElev":
				elev = append(elev, elevator{0, 0, "Idle"})

			case message.Type == "elevOffline":
				elev = append(elev[:message.Value], elev[message.Value+1:]...)

			case message.Type == "newOrder":
				switch{
				case message.Content == "inside":
					inside[message.Floor] = 1
					message.Type = "writeInside"
					message.Value = 1
					fileOutChan <- message
				case message.Content == "outsideUp":
					outUp[message.Floor] = 1
				case message.Content == "outsideDown":
					outDown[message.Floor] = 1
				}
				message.Type = "signal"
				message.Value = 1
				commanderChan <- message

				if message.To == 1 {
					bestValue := FLOOR_COUNT
					bestElev := 0
					for i := 1; i < len(elev); i++ {
						if message.Floor == FLOOR_COUNT && elev[i].floorNum == FLOOR_COUNT-1 && elev[i].state == "MovingUp"{
							bestElev = 0
							break
						} else if message.Floor == 1 && elev[i].floorNum == 2 && elev[i].state == "MovingDown"{
							bestElev = 0
							break
						} else if message.Floor - elev[i].floorNum == 1 && elev[i].state == "MovingUp"{
							bestElev = 0
							break
						} else if message.Floor - elev[i].floorNum == -1 && elev[i].state == "MovingDown"{
							bestElev = 0
							break
						}
						if elev[i].state == "Idle"{
							if message.Floor - elev[i].floorNum > 0 && message.Floor - elev[i].floorNum < bestValue{
								bestValue = message.Floor - elev[i].floorNum
								bestElev = i
							} else if elev[i].floorNum - message.Floor > 0 && elev[i].floorNum - message.Floor < bestValue{
								bestValue = elev[i].floorNum - message.Floor
								bestElev = i
							}
						}
					}
					if bestElev != 0 {
						message.To = bestElev
						message.Type = "newTarget"
						commanderChan <- message
					}
				}

			case message.Type == "deleteOrder":
				switch{
				case message.Content == "inside":
					inside[message.Floor] = 0
					message.Type = "writeInside"
					message.Value = 0
					fileOutChan <- message
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

			case message.Type == "newFloor":
				elev[message.From].floorNum = message.Floor

			case message.Type == "floorReached":
				elev[message.From].floorNum = message.Floor
				if inside[message.Floor] == 1 || outUp[message.Floor] == 1 && elev[message.From].state == "MovingUp" || outDown[message.Floor] == 1 && elev[message.From].state == "MovingDown"{
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
						elev[message.From].floorTarget = 0
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
				message.Type = "targetUpdate"
				commanderChan <- message
				if elev[message.To].state == "Idle" {
					message.Type = "command"
					if message.Floor > elev[message.To].floorNum {
						message.Content = "up"
					} else if message.Floor < elev[message.To].floorNum {
						message.Content = "down"
					}
					commanderChan <- message
				}

			case message.Type == "targetUpdate":
				elev[message.From].floorTarget = message.Floor

			case message.Type == "stateUpdate":
				elev[message.From].state = message.Content

				if elev[message.From].state == "Idle"{
					if message.To == 1 {
						for i := 1; i < FLOOR_COUNT + 1; i++ {
							if elev[message.From].floorNum + i < FLOOR_COUNT + 1 && outDown[elev[message.From].floorNum + i] == 1 ||  outUp[elev[message.From].floorNum + i] == 1 || inside[elev[message.From].floorNum + i] == 1 {
								message.To = message.From
								messate.Type = "newTarget"
								message.Floor = elev[message.From].floorNum + i
								commanderChan <- message
								break
							} else if elev[message.From].floorNum - i > 0 && outDown[elev[message.From].floorNum - i] == 1 ||  outUp[elev[message.From].floorNum - i] == 1 || inside[elev[message.From].floorNum - i] == 1 {
								message.To = message.From
								messate.Type = "newTarget"
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