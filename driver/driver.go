package driver

import (
	."fmt"
	."time"
	."../network"
)

const N_FLOORS = 4
const N_BUTTONS = 3

func DriverInit(driverInChan chan Message, driverOutChan chan Message) (bool) {

	floorSensors := []int{SENSOR_FLOOR1, SENSOR_FLOOR2, SENSOR_FLOOR3, SENSOR_FLOOR4}

	buttonChannelMatrix := [][]int{
		{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
		{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
		{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
		{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},}

	if !IOInit() {
		return false
	}
	for floor := 0; floor < N_FLOORS; floor++ {
		if floor != 0 {
			elevSetButtonLamp("outSideDown", floor, 0)
		}
		if floor != N_FLOORS - 1 {
			elevSetButtonLamp("outsideUp", floor, 0)
		}
		elevSetButtonLamp("inside", floor, 0)
	}
	elevSetStopLamp(0)
	elevSetEngineSpeed("stop")
	elevSetDoorOpenLamp(0)
	elevSetFloorIndicator(1)

	inFloor := 0
	for i := 0; i < N_FLOORS; i++ {
		if IOReadBit(floorSensors[i]) != 0 {
			inFloor = 1
		}
	}
	if inFloor == 0 {
		elevSetEngineSpeed("down")
	}
	go driverReader(driverInChan, floorSensors, buttonChannelMatrix)
	go driverWriter(driverOutChan, floorSensors)

	return true
}

func driverReader(driverInChan chan Message, floorSensors[] int, buttonChannelMatrix[][] int) {
	
	buttonSignalLastCheckMatrix := [][]int{{0,0,0},{0,0,0},{0,0,0},{0,0,0}}
	floorSignalLastCheck := []int{0,0,0,0}
	obstrSignalLastCheck := 0
	stopSignalLastCheck := 0
	message := Message{}

	for {
		for i := 0; i < N_FLOORS; i++ {
			for j := 0; j < N_BUTTONS; j++ {
				if IOReadBit(buttonChannelMatrix[i][j]) != buttonSignalLastCheckMatrix[i][j] {

					if buttonSignalLastCheckMatrix[i][j] == 0 {
						switch {
						case j == 2:
							message.Type = "inside"
						case j == 0:
							message.Type = "outsideUp"
						case j == 1:
							message.Type = "outsideDown"
						}
						message.Value = i
						driverInChan <- message
						buttonSignalLastCheckMatrix[i][j] = 1

					} else if buttonSignalLastCheckMatrix[i][j] == 1 {
						buttonSignalLastCheckMatrix[i][j] = 0
					}
				}
			}
		}
		if IOReadBit(STOP) != stopSignalLastCheck {
			if stopSignalLastCheck == 0 {
				message.Type = "stop"
				driverInChan <- message
				stopSignalLastCheck = 1
			} else if stopSignalLastCheck == 1 {
				stopSignalLastCheck = 0
			}
		}
		if IOReadBit(OBSTRUCTION) != obstrSignalLastCheck {
			if obstrSignalLastCheck == 0 {
				message.Type = "obstr"
				driverInChan <- message
				obstrSignalLastCheck = 1
			} else if obstrSignalLastCheck == 1 {
				obstrSignalLastCheck = 0
			}
		}

		for i := 0; i < N_FLOORS; i++ {
			if IOReadBit(floorSensors[i]) != floorSignalLastCheck[i] {
				if floorSignalLastCheck[i] == 0 {
					message.Type = "floorReached"
					message.Floor = i+1
					driverInChan <- message
					floorSignalLastCheck[i] = 1
				} else if floorSignalLastCheck[i] == 1 {
					floorSignalLastCheck[i] = 0
				}
			}
		}
		Sleep(25 * Millisecond)
	}
}

func driverWriter(driverOutChan chan Message, floorSensors[] int) {
	for {
		select {
		case message := <- driverOutChan:
			Println(message)
			switch {
				case message.Type == "engine":
					elevSetEngineSpeed(message.Content)		// 0 = stop, 1 = up, -1 = down
					if message.Content == "stop" {
						for i := 0; i < N_FLOORS; i++ {
							if IOReadBit(floorSensors[i]) != 0 {
								elevSetDoorOpenLamp(1)
							}
						}
					}
				case message.Type == "floorReached":
					elevSetFloorIndicator(message.Floor)
				case message.Type == "inside" || message.Type == "outsideUp" || message.Type == "outsideDown":
					elevSetButtonLamp(message.Content, message.Floor, message.Value)
				case message.Type == "stop":
					elevSetStopLamp(message.Value)
				case message.Type == "door":
					Println("door value", message.Value)
					elevSetDoorOpenLamp(message.Value)
			}
		}
	}
}

func elevSetEngineSpeed(direction string) {
	switch {
	case direction == "stop":
		IOWriteAnalog(MOTOR, 0)
	case direction == "up":
		IOWriteAnalog(MOTORDIR, 0)
		IOWriteAnalog(MOTOR, 2800)
	case direction == "down":
		IOWriteAnalog(MOTORDIR, 1)
		IOWriteAnalog(MOTOR, 2800)
	}
}

func elevSetStopLamp(value int) {
	switch {
	case value == 1:
		IOSetBit(LIGHT_STOP)
	case value == 0:
		IOClearBit(LIGHT_STOP)
	}
}

func elevSetDoorOpenLamp(value int) {
	switch {
	case value == 1:
		IOSetBit(LIGHT_DOOR_OPEN)
	case value == 0:
		IOClearBit(LIGHT_DOOR_OPEN)
	}
}

func elevSetFloorIndicator(floorNum int) {
	switch {
	case floorNum == 1:
		IOClearBit(LIGHT_FLOOR_IND1)
		IOClearBit(LIGHT_FLOOR_IND2)
	case floorNum == 2:
		IOSetBit(LIGHT_FLOOR_IND1)
		IOClearBit(LIGHT_FLOOR_IND2)
	case floorNum == 3:
		IOClearBit(LIGHT_FLOOR_IND1)
		IOSetBit(LIGHT_FLOOR_IND2)
	case floorNum == 4:
		IOSetBit(LIGHT_FLOOR_IND1)
		IOSetBit(LIGHT_FLOOR_IND2)
	}
}

func elevSetButtonLamp(button string, floor int, value int) {
	switch {
	case value == 1:
		switch {
		case button == "inside":
			switch {
			case floor == 1:
				IOSetBit(LIGHT_COMMAND1)
			case floor == 2:
				IOSetBit(LIGHT_COMMAND2)
			case floor == 3:
				IOSetBit(LIGHT_COMMAND3)
			case floor == 4:
				IOSetBit(LIGHT_COMMAND4)
			}
		case button == "outsideUp":
			switch {
			case floor == 1:
				IOSetBit(LIGHT_UP1)
			case floor == 2:
				IOSetBit(LIGHT_UP2)
			case floor == 3:
				IOSetBit(LIGHT_UP3)
			}
		case button == "outsideDown":
			switch {
			case floor == 2:
				IOSetBit(LIGHT_DOWN2)
			case floor == 3:
				IOSetBit(LIGHT_DOWN3)
			case floor == 4:
				IOSetBit(LIGHT_DOWN4)
			}
		}
	case value == 0:
		switch {
		case button == "inside":
			switch {
			case floor == 1:
				IOClearBit(LIGHT_COMMAND1)
			case floor == 2:
				IOClearBit(LIGHT_COMMAND2)
			case floor == 3:
				IOClearBit(LIGHT_COMMAND3)
			case floor == 4:
				IOClearBit(LIGHT_COMMAND4)
			}
		case button == "outsideUp":
			switch {
			case floor == 1:
				IOClearBit(LIGHT_UP1)
			case floor == 2:
				IOClearBit(LIGHT_UP2)
			case floor == 3:
				IOClearBit(LIGHT_UP3)
			}
		case button == "outsideDown":
			switch {
			case floor == 2:
				IOClearBit(LIGHT_DOWN2)
			case floor == 3:
				IOClearBit(LIGHT_DOWN3)
			case floor == 4:
				IOClearBit(LIGHT_DOWN4)
			}
		}
	}
}
