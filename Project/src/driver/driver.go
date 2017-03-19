package driver

/*
#cgo CFLAGS: -std=gnu11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import (
	"C"
)

import (
	"../misc"
	"../def"
)

type Type int

const (
	TypeComedi     Type = 0
	TypeSimulation Type = 1
)

const (
	NumFloors  = int(C.N_FLOORS)
	NumButtons = int(C.N_BUTTONS)
)



func Init(t Type, simulator string) {
	C.elev_init(C.elev_type(t), C.CString(simulator))
}

func SetMotorDirection(dirn def.MotorDirection) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func SetButtonLamp(button def.ButtonType, floor int, value bool) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(misc.B2I(value)))
}

func SetFloorIndicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func SetDoorOpenLamp(value bool) {
	C.elev_set_door_open_lamp(C.int(misc.B2I(value)))
}

func SetStopLamp(value bool) {
	C.elev_set_stop_lamp(C.int(misc.B2I(value)))
}

func GetButtonSignal(button def.ButtonType, floor int) bool {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor))) != 0
}

func GetFloorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func GetStopSignal() bool {
	return int(C.elev_get_stop_signal()) != 0
}

func GetObstructionSignal() bool {
	return int(C.elev_get_obstruction_signal()) != 0
}