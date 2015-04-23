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

	fileInChan := make(chan Message)
	fileOutChan := make(chan Message)
	mainWaitChan := make(chan Message)
	networkSend := make(chan Message)
	networkReceive := make(chan Message)
	commanderChan := make(chan Message)
	aliveChan := make(chan Message)
	timerChan := make(chan Message)
	tickerChan := make(chan Message)
	timeOutChan := make(chan Message)
	driverInChan := make(chan Message)
	driverOutChan := make(chan Message)

	if !DriverInit(driverInChan, driverOutChan){
		Println("Driver init failed!")
		return
	}
	go FileManager(fileInChan, fileOutChan)
	go Timekeeper(tickerChan, timerChan, timeOutChan)
	go Network(networkReceive, networkSend, fileInChan, fileOutChan)
	go LiftState(networkReceive, commanderChan, aliveChan, fileInChan, fileOutChan)
	go Commander(networkSend, commanderChan, aliveChan, tickerChan, timerChan, timeOutChan, driverInChan, driverOutChan)

	select{
		case <- mainWaitChan:
	}
}


/*

-----------------------------           TO DO           -------------------------------------

Heisen går noen ganger helt feil retning enn det den skal

Dørlyset virker ikke


NB! Mulig deadlock i alive-broadcast init

NB! Mulig deadlock / endless go routine spawn i elevOffline network send

Message
Type, Content, Floor, Value, To, From 

Type: "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", 
      "stateUpdate", "offline", "command", "floorReached", "signal"

*/
