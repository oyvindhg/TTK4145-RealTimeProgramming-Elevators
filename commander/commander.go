package commander

import (
	"net"
	."fmt"
	."time"
	."../timer"
	."../network"
	."../liftState"
	."../driver"
)

func InitCommander(networkSend chan Message, commanderChan chan Message, aliveChan chan Message, signalChan chan Message, tickerChan chan string, timerChan chan TimerInput, timeOutChan chan string, driverInChan chan DriverSignal, driverOutChan chan DriverSignal, requestChan chan Request, replyChan chan Reply, MASTER_INIT_IP string, PORT string, FLOOR_COUNT int, ELEV_COUNT int) {
	ownIP := ""
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		Println("Address error: ", err)
	}

	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ownIP = ipnet.IP.String()
			}
		}
	}
	if ownIP == MASTER_INIT_IP  {
		go master(networkSend, commanderChan, aliveChan, signalChan, tickerChan, timerChan, timeOutChan, driverInChan, driverOutChan, requestChan, replyChan)
	}
	go commander(commanderChan, aliveChan, signalChan, tickerChan, timerChan, timeOutChan, driverInChan, driverOutChan, requestChan, replyChan)
	timerChan <- TimerInput{150, Millisecond, "alive"}
}

func commander(commanderChan chan Message, aliveChan chan Message, signalChan chan Message, tickerChan chan string, timerChan chan TimerInput, timeOutChan chan string, driverInChan chan DriverSignal, driverOutChan chan DriverSignal, requestChan chan Request, replyChan chan Reply) {
	notAliveCount := 0
	for {
		select { 							// ADD FLOORREACHED CASE
		case <- tickerChan:
			if notAliveCount == 10 {
				Println("Master dead!")		// IMPLEMENT PANIC
			}
			notAliveCount++
			
		case <- aliveChan:
			notAliveCount = 0

		case command := <- commanderChan:
			switch {
				case command.Content == "command":
					if command.Command == "up" {
						driverOutChan <- DriverSignal{"engine", 0, 1}
						Println("going up")
					} else if command.Command == "down" {
						driverOutChan <- DriverSignal{"engine", 0, -1}
					} else if command.Command == "stop" {
						driverOutChan <- DriverSignal{"engine", 0, 0}
					}
				case command.Content == "taskDone":
					Println("taskDone")
			}
					

		case signal := <- signalChan:
			Println(signal)
			driverOutChan <- DriverSignal{signal.ButtonType, signal.FloorNumber, 0}

		case timeOut := <- timeOutChan:
			Println(timeOut)

		case driverIn := <- driverInChan:
			Println(driverIn)
		}
	}
}

func master(networkSend chan Message, commanderChan chan Message, aliveChan chan Message, signalChan chan Message, tickerChan chan string, timerChan chan TimerInput, timeOutChan chan string, driverInChan chan DriverSignal, driverOutChan chan DriverSignal, requestChan chan Request, replyChan chan Reply) {
	go aliveBroadcast(networkSend, tickerChan, requestChan, replyChan)
	Println("Master")
}

func aliveBroadcast(networkSend chan Message, tickerChan chan string, requestChan chan Request, replyChan chan Reply) {
	Sleep(100 * Millisecond)		// Give enough time for the other elevators to connect
	requestChan <- Request{"elevCount", 0}
	reply := <- replyChan
	elevCount := reply.Number
	computerIDs := make([]string, elevCount + 1)
	for i := 1; i < elevCount; i++ {
		requestChan <- Request{"computerID", i}
		reply = <- replyChan
		computerIDs[i] = reply.Answer
	}
	for {
		for j := 1; j < elevCount; j++ {
			networkSend <- Message{computerIDs[j], "", "imAlive", "", 1, true, 1, 1, "", ""}
		}
		Sleep(100 * Millisecond)
	}
}

// Content = "imAlive", "newElev", "newOrder", "deleteOrder", "newTarget", rankChange",
//           "stateUpdate", "connectionChange", "command", "taskDone", "floorReached"

// RecipientID, SenderID, Content, Command, ElevNumber,
// Online, Rank, FloorNumber, ButtonType, State

// computerID, onlineStatus, rank, floorNum, floorTarget, state, inElev[]

/*
type DriverSignal struct{
	SignalType string  // engine, floorReached, inside, outsideUp, outsideDown, stop, obstr
	FloorNumber int
	Value int
}*/
