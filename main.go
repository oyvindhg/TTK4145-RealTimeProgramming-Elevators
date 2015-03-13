package main

import(
	."fmt"
	."time"
	."./timer"
	."./network"
	."./liftState"
	."./commander"
)

// Content = "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", rankChange",
//           "stateUpdate", "connectionChange", "command", "taskDone", "floorReached"

// RecipientID, SenderID, Content, Command, ElevNumber,
// Online, Rank, FloorNumber, ButtonType, State

// computerID, onlineStatus, rank, floorNum, floorTarget, state, inElev[]

const ELEV_COUNT = 3
const FLOOR_COUNT = 4
const MASTER_INIT_IP = "129.241.187.147"
const PORT = ":20007"

func main(){

	networkReceive := make(chan Message)
	networkSend := make(chan Message)
	commanderChan := make(chan Message)
	aliveChan := make(chan Message)
	signalChan := make(chan Message)
	timerChan := make(chan TimerInput)
	timeOutChan := make(chan string)
	driverInChan := make(chan DriverSignal)
	driverOutChan := make(chan DriverSignal)
	fetchChan := make(chan Message)

	//go InitTimer(timerChan, timeOutChan)
	go InitNetwork(PORT, networkReceive, networkSend)
	go InitLiftState(networkReceive, commanderChan, aliveChan, signalChan, MASTER_INIT_IP, PORT, FLOOR_COUNT, ELEV_COUNT)
	go InitCommander(commanderChan, aliveChan, signalChan, timerChan, timeOutChan, driverInChan, driverOutChan, fetchChan, MASTER_INIT_IP, PORT, FLOOR_COUNT, ELEV_COUNT)

	go sendStuff(networkSend)

	// INSERT ELEGANT SOLUTION FOR STOP BUTTON TERMINATE

	for {
		select{
			case driver := <- driverOutChan:
				Println(driver.SignalType, driver.FloorNumber)
			case timerOut := <- timeOutChan:
				Println(timerOut)
		}
	}
}

func sendStuff(networkSend chan Message){
	
	initElev1 := Message{MASTER_INIT_IP+ PORT, "", "newID", "", 0, false, 0, 0, "", ""}
	message := Message{PORT, "", "newOrder", "", 1, false, 0, 2, "inside", ""}

	Sleep(1*Second)
	networkSend <- initElev1
	for i := 0; i < 2; i++ {
		Sleep(1*Second)
		networkSend <- message
	}
}
