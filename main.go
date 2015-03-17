package main

import(
	."fmt"
	//."time"
	."./timer"
	."./driver"
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
const MASTER_INIT_IP = "192.168.1.157"
const PORT = ":20007"

func main(){

	networkReceive := make(chan Message)
	networkSend := make(chan Message)
	commanderChan := make(chan Message)
	aliveChan := make(chan Message)
	signalChan := make(chan Message)
	tickerChan := make(chan string)
	timerChan := make(chan TimerInput)
	timeOutChan := make(chan string)
	driverInChan := make(chan DriverSignal)
	driverOutChan := make(chan DriverSignal)
	requestChan := make(chan Request)
	replyChan := make(chan Reply)

	if !DriverInit(driverInChan, driverOutChan) {
		Println("Driver init failed!")
		break
	}
	go InitTimer(tickerChan, timerChan, timeOutChan)
	go InitNetwork(PORT, networkReceive, networkSend)
	go InitLiftState(networkReceive, commanderChan, aliveChan, signalChan, requestChan, replyChan, MASTER_INIT_IP, PORT, FLOOR_COUNT, ELEV_COUNT)
	go InitCommander(networkSend, commanderChan, aliveChan, signalChan, tickerChan, timerChan, timeOutChan, driverInChan, driverOutChan, requestChan, replyChan, MASTER_INIT_IP, PORT, FLOOR_COUNT, ELEV_COUNT)

	//go sendStuff(networkSend)

	// INSERT ELEGANT SOLUTION FOR STOP BUTTON TERMINATE

	// FIX POINTER ELEV IN LIFTSTATE

	for {
		select{
			case driver := <- driverOutChan:
				Println(driver.SignalType, driver.FloorNumber)
		}
	}
}

func sendStuff(networkSend chan Message){
	
	initElev1 := Message{MASTER_INIT_IP+ PORT, "", "newID", "", 0, false, 0, 0, "", ""}
	//message := Message{PORT, "", "imAlive", "", 1, false, 1, 1, "", ""}

	//Sleep(1*Millisecond)
	networkSend <- initElev1
	/*
	for i := 0; i < 1; i++ {
		Sleep(1*Second)
		networkSend <- message
	}*/
}
