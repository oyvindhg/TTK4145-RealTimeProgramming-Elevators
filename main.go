package main

import(
	."fmt"
	."time"
	."./timer"
	."./driver"
	."./network"
	."./liftState"
	."./commander"
)

// TO DO

// MÅ ORDNE ENKEL MÅTE FOR COMMANDER Å SENDE TIL NÅVÆRENDE MASTER

// MÅ ORDNE ENKEL MÅTE FOR COMMANDER Å VITE HVILKEN HEIS SOM FÅR DRIVERINPUT

// Kostfunksjon i LiftState  (FloorReached i Commander blir handla av en stateUpdate i Liftstate)

// Death panic

// Master sender heisinfo når ny heis kobler seg på til de andre heisene + velkomstpakke til ny heis

// Driverinputs til Commander

// Driveroutputs til Commander	 Fixed, liftState må ordne når driverOutputs skal sendes til Commander

// Dørtimer i Commander     ---  La til elevSetDoorLamp(1) i driver.go i en if IOReadBit(floorSensors[i]) != 0
//							---  La til ElevNumber i TimerInput structen for å kunne gi beskjed om hvilken heis som åpner/lukker døra og får Open/Idle state

// Task done i commander 	---  Må legge til algoritme i liftState for å sende til alle heisene, kan ikke gjøre dette i network pga f eks aliveBroadcast (ikke til alle)
//							---  La til taskDone som eget case i liftState, liftState skal jo slette ordre også

func main(){


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

//_______________________________________________________________________________________________________________________________


	//networkSend <- Message{MASTER_INIT_IP+ PORT, "129.241.187.148"+ PORT, "newID", "", 0, false, 0, 0, "", ""}

	for {
		select{
			case <- driverInChan:
				//Println(driver.SignalType, driver.FloorNumber + 1)
		}
	}
}


// Content = "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", rankChange",
//           "stateUpdate", "connectionChange", "command", "taskDone", "floorReached", "signal"
// Command = "up", "down", "stop"

// RecipientID, SenderID, Content, Command, ElevNumber,
// Online, Rank, FloorNumber, ButtonType, State

// computerID, onlineStatus, rank, floorNum, floorTarget, state, inElev[]