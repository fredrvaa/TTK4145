package elevator

import (
	. "../def"
)

type StateMachine struct {
	DriverEvents      DriverElevatorEvents
	OrdersEvents      ElevatorOrdersEvents
	DoorTimerResetCh  chan bool
	ErrorTimerResetCh chan bool

	Elevator Elevator
}

func NewStateMachine(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents, doorTimerResetCh chan bool, errorTimerResetCh chan bool) *StateMachine {
	sm := new(StateMachine)
	sm.DriverEvents = driverEvents
	sm.OrdersEvents = ordersEvents
	sm.DoorTimerResetCh = doorTimerResetCh
	sm.ErrorTimerResetCh = errorTimerResetCh
	sm.Elevator.State = ElevatorState{Active: false, Floor: -1, Direction: DirnStop, Behaviour: ElevatorBehaviourIdle}
	return sm
}

func (this *StateMachine) OnInit() {
	this.Elevator.State.Direction = DirnUp
	this.Elevator.State.Behaviour = ElevatorBehaviourMoving
	this.DriverEvents.MotorDirection <- DirnUp
	this.OrdersEvents.State <- this.Elevator.State
	this.ErrorTimerResetCh <- true
}

func (this *StateMachine) OnButtonPressed(button Button) {
	this.OrdersEvents.Order <- Order{button.Floor, OrderType(button.Type), true}
}

func (this *StateMachine) OnButtonReleased(button Button) {}

func (this *StateMachine) OnStopPressed() {
	if this.Elevator.State.Active {
		this.Elevator.State.Active = false
		this.DriverEvents.MotorDirection <- DirnStop
		this.DriverEvents.Light <- LightEvent{LightTypeStop, 0, true}
		this.OrdersEvents.State <- this.Elevator.State
	}
}

func (this *StateMachine) OnStopReleased() {}

func (this *StateMachine) OnFloorReached(floor int) {
	this.Elevator.State.Floor = floor
	this.OrdersEvents.State <- this.Elevator.State
	this.DriverEvents.FloorIndicator <- this.Elevator.State.Floor
	this.ErrorTimerResetCh <- true

	switch this.Elevator.State.Behaviour {
	case ElevatorBehaviourMoving:
		if !this.Elevator.State.Active {
			this.Elevator.State.Active = true
			this.OrdersEvents.State <- this.Elevator.State
		}
		if ShouldStop(this.Elevator) {
			if IsOrderAtFloor(this.Elevator) {
				if this.Elevator.State.Direction == DirnUp {
					this.OrdersEvents.Order <- Order{this.Elevator.State.Floor, OrderCallUp, false}
				} else if this.Elevator.State.Direction == DirnDown {
					this.OrdersEvents.Order <- Order{this.Elevator.State.Floor, OrderCallDown, false}
				}
				this.OrdersEvents.Order <- Order{this.Elevator.State.Floor, OrderCallCommand, false}

				this.DriverEvents.DoorOpen <- true
				this.DoorTimerResetCh <- true

				this.Elevator.State.Behaviour = ElevatorBehaviourDoorOpen
			} else {
				this.Elevator.State.Behaviour = ElevatorBehaviourIdle
				this.Elevator.State.Direction = DirnStop
			}
			this.DriverEvents.MotorDirection <- DirnStop
			this.OrdersEvents.State <- this.Elevator.State
		}
	
	}
}

func (this *StateMachine) OnLocalOrdersUpdated(localOrders Orders) {
	this.Elevator.Orders = localOrders
	for f := 0; f < NumFloors; f++ {
		this.DriverEvents.Light <- LightEvent{LightType(OrderCallCommand), f, this.Elevator.Orders[f][OrderCallCommand]}
	}
	if this.Elevator.State.Active {
		switch this.Elevator.State.Behaviour {
		case ElevatorBehaviourDoorOpen:
			if IsOrderAtFloor(this.Elevator) {
				this.OrdersEvents.Order <- Order{this.Elevator.State.Floor, OrderCallUp, false}
				this.OrdersEvents.Order <- Order{this.Elevator.State.Floor, OrderCallDown, false}
				this.OrdersEvents.Order <- Order{this.Elevator.State.Floor, OrderCallCommand, false}
				this.DoorTimerResetCh <- true
			}
		case ElevatorBehaviourIdle:
			if IsOrderAtFloor(this.Elevator) {
				this.OrdersEvents.Order <- Order{this.Elevator.State.Floor, OrderCallUp, false}
				this.OrdersEvents.Order <- Order{this.Elevator.State.Floor, OrderCallDown, false}
				this.OrdersEvents.Order <- Order{this.Elevator.State.Floor, OrderCallCommand, false}
				this.DoorTimerResetCh <- true
				this.DriverEvents.DoorOpen <- true
				this.Elevator.State.Behaviour = ElevatorBehaviourDoorOpen
				this.OrdersEvents.State <- this.Elevator.State

			} else {
				this.Elevator.State.Direction = GetDirection(this.Elevator)
				if this.Elevator.State.Direction == DirnStop {
					this.Elevator.State.Behaviour = ElevatorBehaviourIdle
				} else {
					this.Elevator.State.Behaviour = ElevatorBehaviourMoving
				}
				this.OrdersEvents.State <- this.Elevator.State
				this.DriverEvents.MotorDirection <- this.Elevator.State.Direction
				this.ErrorTimerResetCh <- true
			}
		}
	}
}

func (this *StateMachine) OnGlobalOrdersUpdated(globalOrders Orders) {
	for f := 0; f < NumFloors; f++ {
		this.DriverEvents.Light <- LightEvent{LightType(OrderCallDown), f, globalOrders[f][OrderCallDown]}
		this.DriverEvents.Light <- LightEvent{LightType(OrderCallUp), f, globalOrders[f][OrderCallUp]}
	}
}

func (this *StateMachine) OnDoorTimerTimeout() {
	switch this.Elevator.State.Behaviour {
	case ElevatorBehaviourDoorOpen:
		if this.Elevator.State.Active {
			this.Elevator.State.Direction = GetDirection(this.Elevator)
			if this.Elevator.State.Direction == DirnStop {
				this.Elevator.State.Behaviour = ElevatorBehaviourIdle
			} else {
				this.Elevator.State.Behaviour = ElevatorBehaviourMoving
			}
			this.OrdersEvents.State <- this.Elevator.State
			this.DriverEvents.MotorDirection <- this.Elevator.State.Direction
			this.DriverEvents.DoorOpen <- false
			this.ErrorTimerResetCh <- true
		}
	}
}

func (this *StateMachine) OnErrorTimerTimeout() {
	switch this.Elevator.State.Behaviour {
	case ElevatorBehaviourMoving:
		this.Elevator.State.Active = false
		this.DriverEvents.MotorDirection <- DirnStop
		this.DriverEvents.Light <- LightEvent{LightTypeStop, 0, true}
		this.OrdersEvents.State <- this.Elevator.State
	}
}
