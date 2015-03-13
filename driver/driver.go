package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go

// Number of signals and lamps on a per-floor basis (excl sensor)
const N_BUTTONS = 3

static const int lampChannelMatrix[N_FLOORS][N_BUTTONS] = {
{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

static const int buttonChannelMatrix[N_FLOORS][N_BUTTONS] = {
{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func driverInit() {

	i int

	// Init hardware
	if (!ioInit())
		return 0;
	
	// Zero all floor button lamps
	for (i = 0; i < N_FLOORS; ++i) {
		if (i != 0) {
			elevSetButtonLamp(BUTTON_CALL_DOWN, i, 0);
		if (i != N_FLOORS - 1)
			elevSetButtonLamp(BUTTON_CALL_UP, i, 0);
			elevSetButtonLamp(BUTTON_COMMAND, i, 0);
	}
	
	// Clear stop lamp, door open lamp, and set floor indicator to ground floor.
	elevSetStopLamp(0);
	elevSetDoorOpenLamp(0);
	elevSetFloorIndicator(0);
	
	// Return success.
	return 1;
}

func Driver(){
	for{
		select{
			case
		}
	}
}
