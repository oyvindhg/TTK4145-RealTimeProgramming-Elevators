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
						
					if message.Content == "inside" && message.From != message.To {
						break
					}
					message.Type = "signal"
					message.Value = 1
					commanderChan <- message

					// Kjør kostfunksjon (hvis noen er idle)

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
					elev[message.From].floorTarget = message.Floor

				case message.Type == "stateUpdate":
					elev[message.From].state = message.Content
					//If State == "Idle": Kjør kostfunksjon (hvis flere bestillinger)
				}
		}
	}
}