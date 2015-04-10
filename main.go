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

// Kostfunksjon i LiftState

// Death panic

// Master sende heisinfo når ny heis kobler seg på til de andre heisene + velkomstpakke til ny heis

// Floor reached case i Commander

// Driverinputs til Commander

// Driveroutputs til Commander

// Dørtimer i Commander

// Task done i Commander

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