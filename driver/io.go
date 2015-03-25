package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go

/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"
import ."fmt"

func IOInit() bool{
	returnValue := int(C.io_init())
	if returnValue == 0 {
		return false
	} else {
		return true
	}
}

func IOSetBit(channel int) {
	C.io_set_bit(C.channel)
}

func IOClearBit(channel int){
	C.io_clear_bit(C.channel)
}

func IOWriteAnalog(channel int, value int){
	C.io_write_analog(C.channel, C.value)
	Println("writing engine")
}

func IOReadBit(channel int) int{
	return int(C.io_read_bit(C.channel))
}

func IOReadAnalog(channel int) int{
	return int(C.io_read_analog(C.channel))
}
