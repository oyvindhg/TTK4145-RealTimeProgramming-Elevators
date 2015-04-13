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

DEATH PANIC

Network send tcp connection fail and elevOffline case

Implement read/write in Liftstate

Implement read/write in Network

AddElev case i Network? Alle sender til Master og så sender den videre kopi til alle?

Kostfunksjon i LiftState  (FloorReached i Commander blir handla av en stateUpdate i Liftstate)

Velkomstpakke til ny heis i Liftstate, inkludert IPliste og ordre

NB! Vi må starte heisene i rekkefølge for at Master faktisk skal bli Master



Fiksa newElev og addElev, mangla bare elevOffline for å være komplett!



Message
Type, Content, Floor, Value, To, From 

Type: "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", 
      "stateUpdate", "offline", "command", "floorReached", "signal"

*/
