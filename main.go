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
	Println("Heia!")
	if !DriverInit(driverOutChan, driverInChan){
		Println("Driver init failed!")
		return
	}
	go FileManager(fileOutChan, fileInChan)
	go Timekeeper(tickerChan, timerChan, timeOutChan)
	go NetworkInit(networkReceive, networkSend, fileOutChan, fileInChan)
	go LiftState(networkReceive, commanderChan, aliveChan, fileOutChan, fileInChan)
	go CommanderInit(networkSend, commanderChan, aliveChan, tickerChan, timerChan, timeOutChan, driverOutChan, driverInChan)

	select{
		case <- mainWaitChan:
	}
}


/*

-----------------------------           TO DO           -------------------------------------

Heisen går noen ganger helt feil retning enn det den skal

Fiks un-constant MASTER_INIT_IP som endres i network om ikke master ip tilgjengelig

Fjern buffere

Legg inn commander case i networkSender så den ikke sender til alle

Legg til sjekk etter 1 sekund i commander

Legg inn en count for antall prøvd å sende eller forskjellige stages

Mangler en funksjonalitet for newOrder kostfunksjon for moving states

						MÅ ORDNES			Folk kan fucke med en stakkar i 4. etasje pga inside orders prioritet

Kanskje legge til en teller? BRUTE FORCE

Under utregning kan states og retning og floorNum få prioritetutdeling


NB! Antar ingen nettverk oppdeling

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
