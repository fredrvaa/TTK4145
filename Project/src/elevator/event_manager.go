package elevator

import (
	. "../def"
	"./timer"
	"time"
)

const interval = 15 * time.Millisecond

func EventManager(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents) {

	doorTimerResetCh := make(chan bool)
	doorTimerTimeoutCh := make(chan bool)
	errorTimerResetCh := make(chan bool)
	errorTimerTimeoutCh := make(chan bool)

	go timer.Timer(3*time.Second, doorTimerResetCh, doorTimerTimeoutCh)
	go timer.Timer(5*time.Second, errorTimerResetCh, errorTimerTimeoutCh)

	//actuator := NewActuator(driverEvents, ordersEvents, doorTimerResetCh, errorTimerResetCh)
	sm := NewStateMachine(driverEvents, ordersEvents, doorTimerResetCh, errorTimerResetCh)	
	sm.OnInit()

	for {
		select {

		case buttonEvent := <-driverEvents.Button:
			if buttonEvent.State {
				sm.OnButtonPressed(buttonEvent.Button)
			} else {
				sm.OnButtonReleased(buttonEvent.Button)
			}

		case stop := <-driverEvents.Stop:
			if stop {
				sm.OnStopPressed()
			} else {
				sm.OnStopReleased()
			}

		case floor := <-driverEvents.Floor:
			sm.OnFloorReached(floor)

		case localOrders := <-ordersEvents.LocalOrders:
			sm.OnLocalOrdersUpdated(localOrders)

		case globalOrders := <-ordersEvents.GlobalOrders:
			sm.OnGlobalOrdersUpdated(globalOrders)

		case <-doorTimerTimeoutCh:
			sm.OnDoorTimerTimeout()

		case <-errorTimerTimeoutCh:
			sm.OnErrorTimerTimeout()

		case <-time.After(interval):
		}
	}
}
