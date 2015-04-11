package main

import(
	."fmt"
	."./timer"
	."./driver"
	."./network"
	."./liftState"
	."./commander"
)

func main(){

	mainWaitChan := make(chan int)
	networkReceive := make(chan Message)
	networkSend := make(chan Message)
	commanderChan := make(chan Message)
	aliveChan := make(chan Message)
	tickerChan := make(chan string)
	timerChan := make(chan TimerInput)
	timeOutChan := make(chan string)
	driverInChan := make(chan DriverSignal)
	driverOutChan := make(chan DriverSignal)

	if !DriverInit(driverInChan, driverOutChan){
		Println("Driver init failed!")
		return
	}
	go Timekeeper(tickerChan, timerChan, timeOutChan)
	go Network(networkReceive, networkSend)
	go LiftState(networkReceive, commanderChan, aliveChan)
	go Commander(networkSend, commanderChan, aliveChan, tickerChan, timerChan, timeOutChan, driverInChan, driverOutChan)

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

Death panic

Kostfunksjon i LiftState  (FloorReached i Commander blir handla av en stateUpdate i Liftstate)

Master sender heisinfo når ny heis kobler seg på til de andre heisene + velkomstpakke til ny heis

FIX APPEND TO SLICE (ELEV)

Driverinputs til Commander

*/