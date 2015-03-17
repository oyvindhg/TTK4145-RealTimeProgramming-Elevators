package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go

import (
	."../commander"
)

const N_BUTTONS = 3
const N_FLOORS = 4

const int buttonChannelMatrix[N_FLOORS][N_BUTTONS] = {
{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

const int floorSensorMatrix[N_FLOORS = {SENSOR_FLOOR1, SENSOR_FLOOR2, SENSOR_FLOOR3, SENSOR_FLOOR4}

func DriverInit(driverInChan chan DriverSignal, driverOutChan chan DriverSignal) (int) {

	if !ioInit() {
		return 0
	}
	elevSetEngineSpeed(0)

	for floor := 0; floor < N_FLOORS; ++floor {
		if foor != 0 {
			elevSetButtonLamp("outSideDown", floor, 0)
		}
		if floor != N_FLOORS - 1 {
			elevSetButtonLamp("outsideUp", floor, 0)
		}
		elevSetButtonLamp("inside", floor, 0)
	}

	elevSetStopLamp(0)
	elevSetDoorOpenLamp(0)
	elevSetFloorIndicator(1)
	
	go driverReader(driverInChan)
	go driverWriter(driverOutChan)

	return 1
}

// DriverSignal: SignalType string, FloorNumber int, Value int
// Type: engine, floorReached, inside, outsideUp, outsideDown, stop, obstr

func driverReader(driverInChan chan DriverSignal) {
	buttonSignalLastCheck := 0
	floorSignalLastCheck := 0
	obstrSignalLastCheck := 0
	stopSignalLastCheck := 0
	floorSignalNum := 0
	for {
		for i := 0; i < N_FLOORS; i++ {
			for j := 0; j < N_BUTTONS; j++ {
				if ioReadBit(buttonChannelMatrix[i][j]) != buttonSignalLastCheck {
					if buttonSignalLastCheck == 0 {
						switch {
						case i == 2:
							driverInChan <- DriverMessage{"inside", j, 1}
						case i == 0:
							driverInChan <- DriverMessage{"outsideUp", j, 1}
						case i == 1:
							driverInChan <- DriverMessage{"outsideDown", j, 1}
						}
						buttonSignalLastChec = 1
					} else if buttonSignalLastCheck == 1 {
						buttonSignalLastCheck = 0
					}
				}
			}
		}
		if ioReadBit(STOP) != stopSignalLastCheck {
			if stopSignalLastCheck == 0 {
				driverInChan <- DriverMessage{"stop", 0, 1}
				stopSignalLastCheck = 1
			} else if stopSignalLastCheck == 1 {
				stopSignalLastCheck = 0
			}
		}
		if ioReadBit(OBSTRUCTION) != obstrSignalLastCheck {
			if obstrSignalLastCheck == 0 {
				driverInChan <- DriverMessage{"obstr", 0, 1}
				obstrSignalLastCheck = 1
			} else if obstrSignalLastCheck == 1 {
				obstrSignalLastCheck = 0
			}
		}
		if ioReadBit(OBSTRUCTION) != obstrSignalLastCheck {
			if obstrSignalLastCheck == 0 {
				driverInChan <- DriverMessage{"obstr", 0, 1}
				obstrSignalLastCheck = 1
			} else if obstrSignalLastCheck == 1 {
				obstrSignalLastCheck = 0
			}
		}
		for i := 0; i < N_FLOORS; i++ {
			if ioReadBit(floorSensorMatrix[i]) != floorSignalLastCheck {
				if floorSignalLastCheck == 0 {
					driverInChan <- DriverMessage{"floorReached", i, 1}
					floorSignalLastCheck = 1
				} else if floorSignalLastCheck == 1 {
					floorSignalLastCheck = 0
				}
			}
		}
	}
}

func driverWriter(driverOutChan chan DriverSignal) {
	for {
		select {
		case command <- driverOutChan:
			switch {
				case command.SignalType == "engine":
					elevSetEngineSpeed(command.Value)		// 0 = stop, 1 = up, -1 = down
				case command.SignalType == "floorReached":
					elevSetFloorIndicator(command.FloorNumber, command.Value)
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
		ioWriteAnalog(MOTOR, 0)
	case value == 1:
		ioWriteAnalog(MOTORDIR, 1)
		ioWriteAnalog(MOTOR, 100)
	case value == -1:
		ioWriteAnalog(MOTORDIR, 0)
		ioWriteAnalog(MOTOR, 100)
	}
}

func elevSetStopLamp(value int) {
	switch {
	case value == 1:
		ioSetBit(LIGHT_STOP)
	case value == 0:
		ioClearBit(LIGHT_STOP)
	}
}

func elevSetDoorOpenLamp(value int) {
	switch {
	case value == 1:
		ioSetBit(LIGHT_DOOR_OPEN)
	case value == 0:
		ioClearBit(LIGHT_DOOR_OPEN)
	}
}

func elevSetFloorIndicator(floorNum int) {
	switch {
	case floorNum == 1:
		ioClearBit(LIGHT_FLOOR_IND1)
		ioClearBit(LIGHT_FLOOR_IND2)
	case floorNum == 2:
		ioSetBit(LIGHT_FLOOR_IND1)
		ioClearBit(LIGHT_FLOOR_IND2)
	case floorNum == 3:
		ioClearBit(LIGHT_FLOOR_IND1)
		ioSetBit(LIGHT_FLOOR_IND2)
	case floorNum == 4:
		ioSetBit(LIGHT_FLOOR_IND1)
		ioSetBit(LIGHT_FLOOR_IND2)
	}
}

func elevSetButtonLamp(button string, floor int, value int) {
	switch {
	case value == 1:
		switch {
		case button == "inside":
			switch {
			case floor == 1:
				ioSetBit(LIGHT_COMMAND1)
			case floor == 2:
				ioSetBit(LIGHT_COMMAND2)
			case floor == 3:
				ioSetBit(LIGHT_COMMAND3)
			case floor == 4:
				ioSetBit(LIGHT_COMMAND4)
			}
		case button == "outsideUp":
			switch {
			case floor == 1:
				ioSetBit(LIGHT_UP1)
			case floor == 2:
				ioSetBit(LIGHT_UP2)
			case floor == 3:
				ioSetBit(LIGHT_UP3)
			}
		case button == "outsideDown":
			switch {
			case floor == 2:
				ioSetBit(LIGHT_DOWN2)
			case floor == 3:
				ioSetBit(LIGHT_DOWN3)
			case floor == 4:
				ioSetBit(LIGHT_DOWN4)
			}
		}
	case value == 0:
		switch {
		case button == "inside":
			switch {
			case floor == 1:
				ioClearBit(LIGHT_COMMAND1)
			case floor == 2:
				ioClearBit(LIGHT_COMMAND2)
			case floor == 3:
				ioClearBit(LIGHT_COMMAND3)
			case floor == 4:
				ioClearBit(LIGHT_COMMAND4)
			}
		case button == "outsideUp":
			switch {
			case floor == 1:
				ioClearBit(LIGHT_UP1)
			case floor == 2:
				ioClearBit(LIGHT_UP2)
			case floor == 3:
				ioClearBit(LIGHT_UP3)
			}
		case button == "outsideDown":
			switch {
			case floor == 2:
				ioClearBit(LIGHT_DOWN2)
			case floor == 3:
				ioClearBit(LIGHT_DOWN3)
			case floor == 4:
				ioClearBit(LIGHT_DOWN4)
			}
		}
	}
}