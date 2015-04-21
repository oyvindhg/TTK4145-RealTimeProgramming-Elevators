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
						
					case message.Content == "outsideUp":
						outUp[message.Floor] = 1
					case message.Content == "outsideDown":
						outDown[message.Floor] = 1
					}
						
					if message.Content == "inside" && message.From != message.To {							//Vil dette egentlig skje?
						break
					}
					message.Type = "signal"
					message.Value = 1
					commanderChan <- message

					if message.To == 1 {								//KOSTFUNKSJON
						bestValue := FLOOR_COUNT
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
						message.To = bestElev
						message.Type = "newTarget"
						commanderChan <- message
					}

				case message.Type == "deleteOrder":
					switch{
					case message.Content == "inside":
						inside[message.Floor] = 0
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
					elev[message.To].floorTarget = message.Floor
					message.Type = "targetUpdate"						//Hvor er denne typen i Driver?
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

					// DETTE ER NYTT

					if elev[message.From].state == "Idle"{
						if message.To == 1 {						// KOSTFUNKSJON
							bestFloor := 0
							i := 0
							for{
								if elev[message.From].floorNum + i <= FLOOR_COUNT + 1{
									if outDown[elev[message.From].floorNum + i] == 1 ||  outUp[elev[message.From].floorNum + i] == 1 {
										bestFloor = elev[message.From].floorNum + i
										message.To = message.From
										message.Type = "newTarget"
										message.Floor = bestFloor
										commanderChan <- message
										break
									}
								} else if elev[message.From].floorNum - i > 0{
									if outDown[elev[message.From].floorNum - i] == 1 ||  outUp[elev[message.From].floorNum - i] == 1 {
										bestFloor = elev[message.From].floorNum - i
										message.To = message.From
										message.Type = "newTarget"
										message.Floor = bestFloor
										commanderChan <- message
										break
									}
								} else{
									break
								}
								i ++
							}
						}
					}

					// HER SLUTTER DET NYE

					if elev[message.To].state == "Idle" {
						message.Type = "command"
						if message.Floor > elev[message.To].floorNum {		//wtf.. message.Floor er vel ikke bestemt i stateUpdate? Er det ikke her master skal sjekke om heisen skal få en ny oppgave?
							message.Content = "up"
						} else if message.Floor < elev[message.To].floorNum {
							message.Content = "down"
						}
						commanderChan <- message
					}
				}
		}
	}
}