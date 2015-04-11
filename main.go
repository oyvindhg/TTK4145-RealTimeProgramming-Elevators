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

	mainWaitChan := make(chan int)
	networkSend := make(chan Message)
	networkReceive := make(chan Message)
	commanderChan := make(chan Message)
	aliveChan := make(chan Message)
	timerChan := make(chan Message)
	tickerChan := make(chan string)
	timeOutChan := make(chan string)
	driverInChan := make(chan Message)
	driverOutChan := make(chan Message)

	if !DriverInit(driverInChan, driverOutChan){
		Println("Driver init failed!")
		return
	}
	go Timekeeper(tickerChan, timerChan, timeOutChan)
	go Network(networkReceive, networkSend)
	go LiftState(networkReceive, commanderChan, aliveChan)
	go Commander(networkSend, commanderChan, aliveChan, tickerChan, timerChan, timeOutChan, driverInChan, driverOutChan)
	go ReadIP()
	select{
		case <- mainWaitChan:
	}
}


/*

Message
Type, Content, Floor, Value, To, From 

Type: "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", rankChange",
          "stateUpdate", "connectionChange", "command", "floorReached", "signal"


-----------------------------           TO DO           -------------------------------------

FIX DEADLOCK IN NETWORK BROADCAST

Death panic

Kostfunksjon i LiftState  (FloorReached i Commander blir handla av en stateUpdate i Liftstate)

Master sender heisinfo når ny heis kobler seg på til de andre heisene + velkomstpakke til ny heis

FIX APPEND TO SLICE (ELEV)

Driverinputs til Commander

*/
