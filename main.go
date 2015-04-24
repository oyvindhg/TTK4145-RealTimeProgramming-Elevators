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
	go NetworkInit(networkReceive, networkSend, fileOutChan, fileInChan)
	go LiftStateInit(networkReceive, commanderChan, aliveChan, fileOutChan, fileInChan)
	go CommanderInit(networkSend, commanderChan, aliveChan, tickerChan, timerChan, timeOutChan, driverOutChan, driverInChan)

	select{
		case <- mainWaitChan:
	}
}


/*

-----------------------------           TO DO           -------------------------------------

Heisen går noen ganger helt feil retning enn det den skal

Heis på vei bort fra etasje - trykker etasjen den var i -> stopper


AMMAGAAAD FIX FLOORUPDATE IN FLOORREACHED GEEZUS

HUSK Å LEGGE INN POINTER OG REFERENCE I NETWORKRECEIVER TIL IPLIST SLICE


NB! Når ordre for en heis i 4. etasje bestilles opp fra 3. og så 2. til tom kø
	vil den ikke kjøre ned til 2. etasje først, men fikse kun 3. etasje og går ut ifra
	at de andre heisene fikser duden i 2. etasje

NB! DoorTimer skriver og leser til en samme global variabel kanskje helt samtidig

NB! Mulig deadlock i alive-broadcast init

NB! Mulig deadlock / endless go routine spawn i elevOffline network send

Message
Type, Content, Floor, Value, To, From 

Type: "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", 
      "stateUpdate", "offline", "command", "floorReached", "signal"

*/
