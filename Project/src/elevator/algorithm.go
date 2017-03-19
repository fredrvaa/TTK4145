package elevator

import (
	. "../def"
)

func IsOrderAtFloor(elev Elevator) bool {
	for o := 0; o < NumTypes; o++ {
		if elev.Orders[elev.State.Floor][o] {
			return true
		}
	}
	return false
}

func GetDirection(elev Elevator) MotorDirection {
	switch elev.State.Direction {
	case DirnUp:
		if checkAbove(elev) {
			return DirnUp
		} else if checkBelow(elev) {
			return DirnDown
		} else {
			return DirnStop
		}
	case DirnStop, DirnDown:
		if checkBelow(elev) {
			return DirnDown
		} else if checkAbove(elev) {
			return DirnUp
		} else {
			return DirnStop
		}
	}
	return DirnStop
}

func ShouldStop(elev Elevator) bool {
	switch elev.State.Direction {
	case DirnDown:
		return elev.Orders[elev.State.Floor][OrderCallCommand] ||
			elev.Orders[elev.State.Floor][OrderCallDown] ||
			!checkBelow(elev)
	case DirnUp:
		return elev.Orders[elev.State.Floor][OrderCallCommand] ||
			elev.Orders[elev.State.Floor][OrderCallUp] ||
			!checkAbove(elev)
	}
	return true
}

/*
func GetOrdersToClear(orders Orders, floor int, direction MotorDirection) {

	if sm.State.Direction == DirnUp {
		sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallUp, false}
	} else if sm.State.Direction == DirnDown {
		sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallDown, false}
	}
	sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallCommand, false}
}*/

func checkAbove(elev Elevator) bool {
	for f := elev.State.Floor + 1; f < NumFloors; f++ {
		for o := 0; o < NumTypes; o++ {
			if elev.Orders[f][o] {
				return true
			}
		}
	}
	return false
}

func checkBelow(elev Elevator) bool {
	for f := 0; f < elev.State.Floor; f++ {
		for o := 0; o < NumTypes; o++ {
			if elev.Orders[f][o] {
				return true
			}
		}
	}
	return false
}

func ClearOrdersAtCurrentFloor(elev Elevator) Elevator {
	elev.Orders[elev.State.Floor][ButtonCallCommand] = false

	switch elev.State.Direction {
	case DirnUp:
		elev.Orders[elev.State.Floor][ButtonCallUp] = false
		if !checkAbove(elev) {
			elev.Orders[elev.State.Floor][ButtonCallDown] = false
		}
	case DirnDown:
		elev.Orders[elev.State.Floor][ButtonCallDown] = false
		if !checkBelow(elev) {
			elev.Orders[elev.State.Floor][ButtonCallUp] = false
		}
	default:
		elev.Orders[elev.State.Floor][ButtonCallUp] = false
		elev.Orders[elev.State.Floor][ButtonCallDown] = false
	}

	return elev
}
