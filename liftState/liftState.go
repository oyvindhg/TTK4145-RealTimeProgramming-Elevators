package liftState

import (
	"net"
	."fmt"
	."time"
	."../network"
)

const ELEV_COUNT = 3
const FLOOR_COUNT = 4
const MASTER_INIT_IP = "129.241.187.148"

type elevator struct {
	computerID string
	onlineStatus bool
	rank int
	floorNum int
	floorTarget int
	state string		//Idle, Open, MovingUp, MovingDown
}

func LiftState(networkReceive chan Message, commanderChan chan Message, aliveChan chan Message) {
	elev := make([]elevator, 1)
	inside 	:= make([][]int, ELEV_COUNT - 1, FLOOR_COUNT - 1)
	outUp 	:= make([]int, FLOOR_COUNT - 1)
	outDown	:= make([]int, FLOOR_COUNT - 1)

	addresses, err := net.InterfaceAddrs()
	if err != nil {
		Println("Address error: ", err)
	}
	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				elev[0].computerID = ipnet.IP.String()+ PORT
			}
		}
	}

	if elev[0].computerID == MASTER_INIT_IP+ PORT {
		go aliveBroadcast(commanderChan, &elev)
	}

	commanderChan <- Message{MASTER_INIT_IP+ PORT, elev[0].computerID, "newID", "", 0, true, 0, 0, "", ""}

	for{
		select{
			case message := <- networkReceive:
				switch{
					case message.Content == "imAlive":
						aliveChan <- message

					case message.Content == "command" || message.Content == "taskDone":
						commanderChan <- message

					case message.Content == "newID": 
						temp := make([]elevator, len(elev) + 1, cap(elev) + 1)

						for i := range elev {
							temp[i] = (elev)[i]
						}
						elev = temp
						(elev)[len(elev) - 1].computerID = message.SenderID
						(elev)[len(elev) - 1].onlineStatus = message.Online
						(elev)[len(elev) - 1].rank = len(elev)
						(elev)[len(elev) - 1].floorNum = 0
						(elev)[len(elev) - 1].floorTarget = 0
						(elev)[len(elev) - 1].state = "Idle"						//SEND NEWID TO ALL ELEVATORS IF MASTER
						
					case message.Content == "connectionChange":
						(elev)[message.ElevNumber].onlineStatus = message.Online

					case message.Content == "rankChange":
						(elev)[message.ElevNumber].rank = message.Rank

					case message.Content == "newOrder":
						switch{
						case message.ButtonType == "inside":
							inside[message.ElevNumber][message.FloorNumber - 1] = 1
							
						case message.ButtonType == "outsideUp":
							outUp[message.FloorNumber] = 1
						case message.ButtonType == "outsideDown":
							outDown[message.FloorNumber] = 1
						}
							//CHECK IF NOT OWN ELEVATOR && INSIDE DON'T SEND SIGNALCHAN
						commanderChan <- message

						// Kjør kostfunksjon (hvis noen er idle)

					case message.Content == "deleteOrder":
						switch{
						case message.ButtonType == "inside":
							inside[message.ElevNumber][message.FloorNumber - 1] = 0
						case message.ButtonType == "outsideUp":
							outUp[message.FloorNumber] = 0
						case message.ButtonType == "outsideDown":
							outDown[message.FloorNumber] = 0
						}
							//CHECK IF NOT OWN ELEVATOR && INSIDE DON'T SEND SIGNALCHAN
						commanderChan <- message

					case message.Content == "newTarget":
						(elev)[message.ElevNumber].floorTarget = message.FloorNumber

					case message.Content == "stateUpdate":
						(elev)[message.ElevNumber].state = message.State
						//If State == "Idle": Kjør kostfunksjon (hvis flere bestillinger)
				}
		}
	}
}

func aliveBroadcast(commanderChan chan Message, elev *[]elevator) {
	Sleep(100 * Millisecond)		// Give enough time for the other elevators to connect
	for {
		for i := 1; i < ELEV_COUNT + 1; i++ {
			if i < len(*elev) {
				commanderChan <- Message{(*elev)[i].computerID, "", "imAlive", "", 0, true, 0, 0, "", ""}
			}
		}
		Sleep(100 * Millisecond)
	}
}


// --- MESSAGE CONTENT ---
// Content = "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", rankChange",
//           "stateUpdate", "connectionChange", "command", "taskDone", "signal"


// --- MESSAGE STRUCT ---
// RecipientID, SenderID, Content, Command, ElevNumber,
// Online, Rank, FloorNumber, ButtonType, State


// --- ELEVATOR STRUCT ---
// computerID, onlineStatus, rank, floorNum, floorTarget, state, inElev[]