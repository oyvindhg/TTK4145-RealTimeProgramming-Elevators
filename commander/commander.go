package commander

import (
	"net"
	."fmt"
	."../timer"
	."../network"
	//."../liftState"
	//."../driver"
)

type DriverSignal struct{
	SignalType string  // engine, floorReached, inside, outsideUp, outsideDown, stop, obstr
	FloorNumber int
	Engine string
}

func InitCommander(commanderChan chan Message, aliveChan chan Message, signalChan chan Message, timerChan chan TimerInput, timeOutChan chan string, driverInChan chan DriverSignal, driverOutChan chan DriverSignal, fetchChan chan Message, MASTER_INIT_IP string, PORT string, FLOOR_COUNT int, ELEV_COUNT int) {
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
		go master(commanderChan, aliveChan, signalChan, timerChan, timeOutChan, driverInChan, driverOutChan, fetchChan)
	}
	go commander(commanderChan, aliveChan, signalChan, timerChan, timeOutChan, driverInChan, driverOutChan, fetchChan)
}

func commander(commanderChan chan Message, aliveChan chan Message, signalChan chan Message, timerChan chan TimerInput, timeOutChan chan string, driverInChan chan DriverSignal, driverOutChan chan DriverSignal, fetcChan chan Message) {
	for {
		select { //LEGG TIL FLOORREACHED CASE
		case alive := <- aliveChan:
			timerInput := TimerInput{25000, "Millisecond", ""}
			timerChan <- timerInput
			Println(alive)

		case command := <- commanderChan:
			Println(command)

		case signal := <- signalChan:
			Println(signal)
			driverOutChan <- DriverSignal{signal.ButtonType, signal.FloorNumber, ""}

		case timeOut := <- timeOutChan:
			Println(timeOut)

		case driverIn := <- driverInChan:
			Println(driverIn)
		}
	}
}

func master(commanderChan chan Message, aliveChan chan Message, signalChan chan Message, timerChan chan TimerInput, timeOutChan chan string, driverInChan chan DriverSignal, driverOutChan chan DriverSignal, fetcChan chan Message) {
	Println("Master")
}