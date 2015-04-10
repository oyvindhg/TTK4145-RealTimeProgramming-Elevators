package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go

//import ."fmt"

const N_BUTTONS = 3
const N_FLOORS = 4

type DriverSignal struct{
	SignalType string  // engine, floorReached, inside, outsideUp, outsideDown, stop, obstr
	FloorNumber int
	Value int
}

func DriverInit(driverInChan chan DriverSignal, driverOutChan chan DriverSignal) (bool) {

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
	elevSetEngineSpeed(0)
	elevSetDoorOpenLamp(0)
	elevSetFloorIndicator(1)
	inFloor := 0
	for i := 0; i < N_FLOORS; i++ {
		if IOReadBit(floorSensors[i]) != 0 {
			inFloor = 1
		}
	}
	if inFloor == 0 {
		elevSetEngineSpeed(-1)
	}
	go driverReader(driverInChan, floorSensors, buttonChannelMatrix)
	go driverWriter(driverOutChan)

	return true
}

func driverReader(driverInChan chan DriverSignal, floorSensors[] int, buttonChannelMatrix[][] int) {
	
	buttonSignalLastCheckMatrix := [][]int{{0,0,0},{0,0,0},{0,0,0},{0,0,0}}
	floorSignalLastCheck := []int{0,0,0,0}
	obstrSignalLastCheck := 0
	stopSignalLastCheck := 0
	for {
		for i := 0; i < N_FLOORS; i++ {
			for j := 0; j < N_BUTTONS; j++ {
				if IOReadBit(buttonChannelMatrix[i][j]) != buttonSignalLastCheckMatrix[i][j] {
					if buttonSignalLastCheckMatrix[i][j] == 0 {
						switch {
						case j == 2:
							driverInChan <- DriverSignal{"inside", i, 1}
						case j == 0:
							driverInChan <- DriverSignal{"outsideUp", i, 1}
						case j == 1:
							driverInChan <- DriverSignal{"outsideDown", i, 1}
						}
						buttonSignalLastCheckMatrix[i][j] = 1
					} else if buttonSignalLastCheckMatrix[i][j] == 1 {
						buttonSignalLastCheckMatrix[i][j] = 0
					}
				}
			}
		}
		if IOReadBit(STOP) != stopSignalLastCheck {
			if stopSignalLastCheck == 0 {
				driverInChan <- DriverSignal{"stop", 0, 1}
				stopSignalLastCheck = 1
			} else if stopSignalLastCheck == 1 {
				stopSignalLastCheck = 0
			}
		}
		if IOReadBit(OBSTRUCTION) != obstrSignalLastCheck {
			if obstrSignalLastCheck == 0 {
				driverInChan <- DriverSignal{"obstr", 0, 1}
				obstrSignalLastCheck = 1
			} else if obstrSignalLastCheck == 1 {
				obstrSignalLastCheck = 0
			}
		}

		for i := 0; i < N_FLOORS; i++ {
			if IOReadBit(floorSensors[i]) != floorSignalLastCheck[i] {
				if floorSignalLastCheck[i] == 0 {
					driverInChan <- DriverSignal{"floorReached", i+1, 1}
					floorSignalLastCheck[i] = 1
				} else if floorSignalLastCheck[i] == 1 {
					floorSignalLastCheck[i] = 0
				}
			}
		}
	}
}

func driverWriter(driverOutChan chan DriverSignal) {
	for {
		select {
		case command := <- driverOutChan:
			switch {
				case command.SignalType == "engine":
					elevSetEngineSpeed(command.Value)		// 0 = stop, 1 = up, -1 = down
				case command.SignalType == "floorReached":
					elevSetFloorIndicator(command.FloorNumber)
				case command.SignalType == "inside" || command.SignalType == "outsideUp" || command.SignalType == "outsideDown":
					elevSetButtonLamp(command.SignalType, command.FloorNumber, command.Value)
				case command.SignalType == "stop":
					elevSetStopLamp(command.Value)
			}
		}
	}
}

func elevSetEngineSpeed(value int) {
	switch {
	case value == 0:
		IOWriteAnalog(MOTOR, 0)
	case value == 1:
		IOWriteAnalog(MOTORDIR, 0)
		IOWriteAnalog(MOTOR, 2800)
	case value == -1:
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
