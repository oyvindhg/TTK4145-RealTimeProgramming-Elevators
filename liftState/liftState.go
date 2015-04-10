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


func liftState(elev *[]elevator, inside[][] int, outUp[] int, outDown[] int, networkReceive chan Message, commanderChan chan Message, aliveChan chan Message, signalChan chan Message, MASTER_INIT_IP string, FLOOR_COUNT int, ELEV_COUNT int) {
	for{
		select{
			case message := <- networkReceive:
				switch{
					case message.Content == "imAlive":
						aliveChan <- message

					case message.Content == "command" || message.Content == "taskDone":
						commanderChan <- message

					case message.Content == "newID": 
						temp := make([]elevator, len(*elev) + 1, cap(*elev) + 1)

						for i := range *elev {
							temp[i] = (*elev)[i]
						}
						*elev = temp
						(*elev)[len(*elev) - 1].computerID = message.SenderID
						(*elev)[len(*elev) - 1].onlineStatus = true
						(*elev)[len(*elev) - 1].rank = len(*elev)   					//NEED BETTER RANK ALGORITHM
						(*elev)[len(*elev) - 1].floorNum = 0
						(*elev)[len(*elev) - 1].floorTarget = 0
						(*elev)[len(*elev) - 1].state = "Idle"
					case message.Content == "connectionChange":
						(*elev)[message.ElevNumber].onlineStatus = message.Online

					case message.Content == "rankChange":
						(*elev)[message.ElevNumber].rank = message.Rank

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
						signalChan <- message

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
						signalChan <- message

					case message.Content == "newTarget":
						(*elev)[message.ElevNumber].floorTarget = message.FloorNumber

					case message.Content == "stateUpdate":
						(*elev)[message.ElevNumber].state = message.State
				}
		}
	}
}

func InitLiftState(networkReceive chan Message, commanderChan chan Message, aliveChan chan Message, signalChan chan Message, requestChan chan Request, replyChan chan Reply, MASTER_INIT_IP string, PORT string, FLOOR_COUNT int, ELEV_COUNT int){
	elev := make([]elevator, 1)
	inside 	:= make([][]int, ELEV_COUNT - 1, FLOOR_COUNT - 1)
	outUp 	:= make([]int, FLOOR_COUNT - 1)
	outDown	:= make([]int, FLOOR_COUNT - 1)
	addresses, err := net.InterfaceAddrs()
	elev[0].onlineStatus = true
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

	go liftState(&elev, inside, outUp, outDown, networkReceive, commanderChan, aliveChan, signalChan, MASTER_INIT_IP, FLOOR_COUNT, ELEV_COUNT)
	go requestHandler(requestChan, replyChan, &elev, outUp, outDown)
}

func requestHandler(requestChan chan Request, replyChan chan Reply, elev *[]elevator, outUp[] int, outDown[] int) {
	reply := Reply{"", 0}
	for {
		select {
		case request := <- requestChan:
			switch {
				case request.Type == "elevCount":
					reply.Number = len(*elev)
					replyChan <- reply
				case request.Type == "computerID":
					reply.Answer = (*elev)[request.ElevNumber].computerID
					replyChan <- reply
			}
		}
	}
}

type elevator struct {
	computerID string
	onlineStatus bool
	rank int
	floorNum int
	floorTarget int
	state string		//Idle, Open, MovingUp, MovingDown
}

type Request struct {
	Type string
	ElevNumber int
}

type Reply struct {
	Answer string
	Number int
}
