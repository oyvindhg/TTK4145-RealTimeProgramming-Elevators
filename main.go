package main

import(
	."fmt"
	."time"
	."./timer"
	."./driver"
	."./network"
	."./liftState"
	."./commander"
	."./fileManager"
)

func main(){

	fileInChan := make(chan Message, 10)
	fileOutChan := make(chan Message, 10)
	mainWaitChan := make(chan Message, 10)
	networkSend := make(chan Message, 10)
	networkReceive := make(chan Message, 10)
	cancelMasterChan := make(chan Message, 10)
	commanderChan := make(chan Message, 10)
	driverOutChan := make(chan Message, 10)
	driverInChan := make(chan Message, 10)
	aliveChan := make(chan Message, 10)
	timerChan := make(chan Message, 10)
	tickerChan := make(chan Message, 10)
	timeOutChan := make(chan Message, 10)
	failureChan := make(chan Message, 10)
	
	if !DriverInit(driverOutChan, driverInChan){
		Println("\n", "Driver init failed!")
		return
	}
	go FileManager(fileOutChan, fileInChan)
	go Timekeeper(tickerChan, timerChan, timeOutChan)
	go NetworkInit(networkReceive, networkSend, fileOutChan, fileInChan, failureChan)
	go LiftState(networkReceive, commanderChan, aliveChan, fileOutChan, fileInChan)
	go CommanderInit(networkSend, commanderChan, aliveChan, tickerChan, timerChan, timeOutChan, driverOutChan, driverInChan, failureChan, cancelMasterChan)
	Println("\n\n\n          --------------------\n          |                  |\n          |   Initializing   |\n          |                  |\n          --------------------\n\n\n")
	Sleep(1050*Millisecond)
	Println("\n\n\n          --------------------\n          |                  |\n          |       DONE       |\n          |                  |\n          --------------------\n\n\n")

	select{
		case <- mainWaitChan:
	}
}
