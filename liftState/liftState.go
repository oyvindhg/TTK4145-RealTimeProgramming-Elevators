package liftState

import (
	."fmt"
	."../network"
)

type elevator struct {
	floorNum int
	floorTarget int
	state string
	isMaster bool
}

func LiftState(networkReceive chan Message, commanderChan chan Message, aliveChan chan Message, fileOutChan chan Message, fileInChan chan Message) {

	message := Message{}
	elev := make([]elevator, ELEV_COUNT + 1)
	inside 	:= make([]int, FLOOR_COUNT+1)
	outUp 	:= make([]int, FLOOR_COUNT+1)
	outDown	:= make([]int, FLOOR_COUNT+1)

	message.Type = "findMaster"
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
	for i:= range elev {
		if i == 0 {
			elev[i].state = "Nil"
		} else {
			elev[i].state = "Uninitialized"
		}
		elev[i].isMaster = false
		elev[i].floorTarget = 0
		elev[i].floorNum = 0
	}

	for {
		select {
		case message = <- networkReceive:
			switch {
			case message.Type == "command" || message.Type == "master" || message.Type == "cancelMaster":
				Println("\n", message.Type, message.Content, "From", message.From, "To", message.To)
				commanderChan <- message

			case message.Type == "masterOffline":
				for i := 1; i < ELEV_COUNT + 1; i++ {
					if elev[i].isMaster {
						message.Type = "masterNumber"
						message.Value = i
						commanderChan <- message
					}
				}

			case message.Type == "broadcast":
				//Println("\n", message.Type, message.Content, "Value =", message.Value, "From", message.From, "To", message.To)
				aliveChan <- message
				







			case message.Type == "newMaster":
				Println("\n", message.Type, message.Content, "Value =", message.Value, "From", message.From, "To", message.To)
				if message.Value == message.To && !elev[message.Value].isMaster {
					message.Type = "master"
					commanderChan <- message
				}
				if message.Value == 0 {
					masterInitCanceled := false
					for i := 1; i < ELEV_COUNT + 1; i++ {
						if elev[i].isMaster {
							masterInitCanceled = true
						}
					}
					if !masterInitCanceled {
						message.Value = message.To
						message.Type = "newMaster"
						commanderChan <- message
					}
				} else {
					for i := 1; i < ELEV_COUNT + 1; i++ {
						if i == message.Value {
							elev[i].isMaster = true
						} else {
							elev[i].isMaster = false
						}
					}
				}







			case message.Type == "findMaster" || message.Type == "elevOnline":
				Println("\n", message.Type, message.Content, "Value =", message.Value, "From", message.From, "To", message.To)
				elev[message.Value].state = "Idle"
				Println("\n", elev)
				if elev[message.To].isMaster {
					for i := 1; i < ELEV_COUNT + 1; i++ {
						if elev[i].state == "Offline" {
							message.Value = i
							break
						}
					}
					message.Type = "newMaster"
					message.To, message.Value = message.Value, message.To
					commanderChan <- message

					message.Type = "addElev"
					message.Value = message.To
					message.To = 0
					commanderChan <- message
					
					message.To = message.Value
					for i := 1; i < message.Value; i++ {
						message.Type = "addElev"
						message.Value = i
						commanderChan <- message
					}
				}






			case message.Type == "addElev":
				Println("\n", message.Type, message.Content, "number", message.Value, "From", message.From, "To", message.To)
				elev[message.Value].state = "Idle"
				Println("\n", elev)
				if message.Value == message.To {
					Println("Sending floorUpdate")
					message.Type = "floorUpdate"
					message.Floor = elev[0].floorNum
					message.To = 0
					commanderChan <- message
				}










			case message.Type == "elevOffline":
				Println("\n", message.Type, message.Content, "Value =", message.Value, "From", message.From, "To", message.To)
				elev[message.Value].state = "Offline"
				Println("\n", elev)
				if elev[message.Value].isMaster {
					nextMaster := 0
					for i := 1; i < ELEV_COUNT + 1; i++ {
						if elev[i].state != "Offline" && elev[i].state != "Uninitialized" {
							nextMaster = i
						}
					}
					message.Value = nextMaster
					message.Type = "newMaster"
					message.To = 0
					commanderChan <- message
				}












































			case message.Type == "newOrder":
				Println("\n", message.Type, message.Content, "Floor =", message.Floor, "From", message.From, "To", message.To)
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
						if elev[i].state == "Idle" {
							if message.Floor - elev[i].floorNum >= 0 && message.Floor - elev[i].floorNum < bestFloor {
								bestFloor = message.Floor - elev[i].floorNum
								bestElev = i
							} else if elev[i].floorNum - message.Floor > 0 && elev[i].floorNum - message.Floor < bestFloor {
								bestFloor = elev[i].floorNum - message.Floor
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

				if message.Content == "inside" && elev[message.From].state == "Idle" {
					message.To = message.From
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
				Println("\n", message.Type, message.Content, "Floor =", message.Floor, "From", message.From, "To", message.To)
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
				Println("\n", message.Type, message.Content, "Floor =", message.Floor, "From", message.From, "To", message.To)
				elev[message.From].floorNum = message.Floor
				if message.From == message.To {
					elev[0].floorNum = message.Floor
				}

			case message.Type == "floorReached":
				Println("\n", message.Type, message.Floor)
				message.Type = "floorUpdate"
				message.To = 0
				commanderChan <- message
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
				Println("\n", message.Type, message.Content, "Floor =", message.Floor, "From", message.From, "To", message.To)
				message.Type = "targetUpdate"
				commanderChan <- message
				message.Type = "command"
				if message.Floor > elev[message.To].floorNum && elev[message.To].state == "Idle" {
					message.Content = "up"
					commanderChan <- message
				} else if message.Floor < elev[message.To].floorNum && elev[message.To].state == "Idle" {
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
				Println("\n", message.Type, message.Content, "Floor =", message.Floor, "From", message.From, "To", message.To)
				elev[message.From].floorTarget = message.Floor

			case message.Type == "stateUpdate":
				Println("\n", message.Type, message.Content, "From", message.From, "To", message.To)
				if message.From == 0 {
					break
				} else {
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
}