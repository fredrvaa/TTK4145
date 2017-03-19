package orders

import (
	. "../def"
	"time"
)

func OrderManager(id string, elevatorEvents ElevatorOrdersEvents, networkEvents OrdersNetworkEvents, guiEvents OrdersGuiEvents) {

	sm := NewStateMachine(id, elevatorEvents, networkEvents, guiEvents)

	for {
		select {
		case order := <-elevatorEvents.Order:
			sm.OnOrderReceived(order)

		case state := <-elevatorEvents.State:
			sm.OnStateReceived(state)

		case newElevator := <-networkEvents.ElevatorNew:
			sm.OnElevatorNew(newElevator)

		case lostElevator := <-networkEvents.ElevatorLost:
			sm.OnElevatorLost(lostElevator)

		case stateEvent := <-networkEvents.RxStateEvent:
			sm.OnStateEventReceived(stateEvent)

		case orderEvent := <-networkEvents.RxOrderEvent:
			sm.OnOrderEventReceived(orderEvent)

		case <-time.After(50 * time.Millisecond):
		}
	}
}