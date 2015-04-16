package liftState

import (
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

func LiftState(networkReceive chan Message, commanderChan chan Message, aliveChan chan Message) {

	message := Message{}
	elev := make([]elevator, 1, ELEV_COUNT + 1)
	inside 	:= make([]int, FLOOR_COUNT+1)
	outUp 	:= make([]int, FLOOR_COUNT+1)
	outDown	:= make([]int, FLOOR_COUNT+1)

	message.Type = "newElev"
	commanderChan <- message

	// READ INSIDEORDERS FROM FILE
	
	for{
		select{
			case message = <- networkReceive:
				switch{
				case message.Type == "broadcast":
					go aliveBroadcast(commanderChan)

				case message.Type == "imAlive":
					aliveChan <- message

				case message.Type == "command":
					commanderChan <- message

				case message.Type == "newElev" || message.Type == "addElev":
					elev = append(elev, elevator{0, 0, "Idle"})

				case message.Type == "elevOffline":
					elev = append(elev[:message.From], elev[message.From+1:]...)

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
					commanderChan <- message

				case message.Type == "newFloor":
					elev[message.From].floorNum = message.Floor

				case message.Type == "floorReached":
					elev[message.From].floorNum = message.Floor
					if inside[message.Floor] == 1 || outUp[message.Floor] == 1 || outDown[message.Floor] == 1 || elev[message.To].floorTarget == 0 {
						message.Type = "command"
						message.Content = "stop"
						commanderChan <- message
						message.Type = "deleteOrder"
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