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


/*

-----------------------------           TO DO           -------------------------------------

Heisen går noen ganger helt feil retning enn det den skal

Floor indicator virker ikke


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
