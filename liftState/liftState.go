package liftState

import (
	"net"
	."fmt"
	."../network"
)

// --- MESSAGE CONTENT ---
// Content = "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", rankChange",
//           "stateUpdate", "connectionChange", "command", "taskDone"


// --- MESSAGE STRUCT ---
// RecipientID, SenderID, Content, Command, ElevNumber,
// Online, Rank, FloorNumber, ButtonType, State


// --- ELEVATOR STRUCT ---
// computerID, onlineStatus, rank, floorNum, floorTarget, state, inElev[]


func liftState(elev[] elevator, outUp[] int, outDown[] int, ownID string, networkReceive chan Message, commanderChan chan Message, aliveChan chan Message, signalChan chan Message, MASTER_INIT_IP string, FLOOR_COUNT int, ELEV_COUNT int) {
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
							temp[i] = elev[i]
						}
						elev = temp
						elev[len(elev) - 1].computerID = message.SenderID
						elev[len(elev) - 1].onlineStatus = true
						elev[len(elev) - 1].rank = len(elev)   					//NEED BETTER RANK ALGORITHM
						elev[len(elev) - 1].floorNum = 0
						elev[len(elev) - 1].floorTarget = 0
						elev[len(elev) - 1].state = "Idle"
						elev[len(elev) - 1].inElev = make([]int, FLOOR_COUNT)		

					case message.Content == "connectionChange":
						elev[message.ElevNumber].onlineStatus = message.Online

					case message.Content == "rankChange":
						elev[message.ElevNumber].rank = message.Rank

					case message.Content == "newOrder":
						switch{
						case message.ButtonType == "inside":
							elev[message.ElevNumber].inElev[message.FloorNumber - 1] = 1
							
						case message.ButtonType == "outsideUp":
							outUp[message.FloorNumber] = 1
						case message.ButtonType == "outsideDown":
							outDown[message.FloorNumber] = 1
						}
							//CHECK IF NOT OWN ELEVATOR && INSIDE DON'T SEND SIGNALCHAN
						signalChan <- message

					case message.Content == "deleteOrder":
						switch{
						case message.ButtonType == "inside":
							elev[message.ElevNumber].inElev[message.FloorNumber - 1] = 0
						case message.ButtonType == "outsideUp":
							outUp[message.FloorNumber] = 0
						case message.ButtonType == "outsideDown":
							outDown[message.FloorNumber] = 0
						}
							//CHECK IF NOT OWN ELEVATOR && INSIDE DON'T SEND SIGNALCHAN
						signalChan <- message

					case message.Content == "newTarget":
						elev[message.ElevNumber].floorTarget = message.FloorNumber

					case message.Content == "stateUpdate":
						elev[message.ElevNumber].state = message.State
				}
		}
	}
}

func InitLiftState(networkReceive chan Message, commanderChan chan Message, aliveChan chan Message, signalChan chan Message, MASTER_INIT_IP string, PORT string, FLOOR_COUNT int, ELEV_COUNT int){
	elev := make([]elevator, 1)
	outUp 	:= make([]int, FLOOR_COUNT - 1)
	outDown	:= make([]int, FLOOR_COUNT - 1)
	ownID := ""
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		Println("Address error: ", err)
	}

	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ownID = ipnet.IP.String()+ PORT
				elev[0].computerID = ownID
			}
		}
	}

	go liftState(elev, outUp, outDown, ownID, networkReceive, commanderChan, aliveChan, signalChan, MASTER_INIT_IP, FLOOR_COUNT, ELEV_COUNT)
}

type elevator struct {
	computerID string
	onlineStatus bool
	rank int
	floorNum int
	floorTarget int
	state string		//Idle, Open, MovingUp, MovingDown
	inElev[] int
}