package main

import(
	."fmt"
	."./timer"
	."./driver"
	."./network"
	."./liftState"
	."./commander"
	."./fileManager"
)

func main(){

	fileInChan := make(chan Message, 0)
	fileOutChan := make(chan Message, 0)
	mainWaitChan := make(chan Message,0)
	networkSend := make(chan Message, 10)
	networkReceive := make(chan Message, 10)
	commanderChan := make(chan Message, 10)
	aliveChan := make(chan Message, 0)
	timerChan := make(chan Message, 0)
	tickerChan := make(chan Message, 0)
	timeOutChan := make(chan Message, 0)
	driverInChan := make(chan Message, 10)
	driverOutChan := make(chan Message, 10)

	if !DriverInit(driverOutChan, driverInChan){
		Println("Driver init failed!")
		return
	}
	go FileManager(fileOutChan, fileInChan)
	go Timekeeper(tickerChan, timerChan, timeOutChan)
	go Network(networkReceive, networkSend, fileOutChan, fileInChan)
	go LiftState(networkReceive, commanderChan, aliveChan, fileOutChan, fileInChan)
	go Commander(networkSend, commanderChan, aliveChan, tickerChan, timerChan, timeOutChan, driverOutChan, driverInChan)

	select{
		case <- mainWaitChan:
	}
}
