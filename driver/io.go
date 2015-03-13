package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go

/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/

import "C"

func ioInit(){
	return int(C.io_init())
}

func ioSetBit(channel int) {
	C.io_set_bit(C.channel)
}

func ioClearBit(channel int){
	C.io_clear_bit(C.channel)
}

func ioWriteAnalog(channel int, value int){
	C.io_write_analog(C.channel, C.value)
}

func ioReadBit(channel int){
	return int(C.io_read_bit(C.channel))
}

func ioReadAnalog(channel int){
	return int(C.io_read_analog(C.channel))
}
