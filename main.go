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

-----------------------------           TO DO           -------------------------------------

FIX DEADLOCK IN NETWORK BROADCAST

Implement read/write in Liftstate

Implement read/write in Network

Kostfunksjon i LiftState  (FloorReached i Commander blir handla av en stateUpdate i Liftstate)



NB! Vi må starte heisene i rekkefølge for at Master faktisk skal bli Master

NB! Mulig deadlock / endless go routine spawn i elevOffline network send

Fiksa newElev, addElev og offlineElev cases og sjekk for tcp



Message
Type, Content, Floor, Value, To, From 

Type: "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", 
      "stateUpdate", "offline", "command", "floorReached", "signal"

*/
