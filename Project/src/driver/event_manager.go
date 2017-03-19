package driver

import (
	"time"
	"../def"
)

func EventManager(elevatorEvents def.DriverElevatorEvents) {


	// Storage of last states to detect change of state
	var lastButtonState [NumFloors][NumButtons]bool
	for f := 0; f < NumFloors; f++ {
		for b := 0; b < NumButtons; b++ {
			lastButtonState[f][b] = false
		}
	}

	lastStopState := false
	lastFloorState := -1//GetFloorSignal()

	var buttonState bool
	var stopState bool
	var floorState int

	for {
		select {
			
		case motorDirection := <-elevatorEvents.MotorDirection:
			SetMotorDirection(motorDirection)

		case light := <-elevatorEvents.Light:
			switch light.LightType {
			case def.LightTypeUp:
				SetButtonLamp(def.ButtonCallUp, light.Floor, light.Value)
			case def.LightTypeDown:
				SetButtonLamp(def.ButtonCallDown, light.Floor, light.Value)
			case def.LightTypeCommand:
				SetButtonLamp(def.ButtonCallCommand, light.Floor, light.Value)
			case def.LightTypeStop:
				SetStopLamp(light.Value)
			}

		case doorOpen := <-elevatorEvents.DoorOpen:
			SetDoorOpenLamp(doorOpen)

		case floorIndicator := <-elevatorEvents.FloorIndicator:
			SetFloorIndicator(floorIndicator)

		case <-time.After(10 * time.Millisecond):
			stopState = GetStopSignal()
			if stopState != lastStopState {
				lastStopState = stopState
				elevatorEvents.Stop <- stopState
			}

			floorState = GetFloorSignal()
			if floorState != lastFloorState {
				lastFloorState = floorState
				if floorState>= 0 && floorState < NumFloors {
					elevatorEvents.Floor <- floorState
				}
			}

			for f := 0; f < NumFloors; f++ {
				for b := 0; b < NumButtons; b++ {
					if def.ButtonType(b) == def.ButtonCallUp && f == NumFloors-1 {
						continue
					}
					if def.ButtonType(b) == def.ButtonCallDown && f == 0 {
						continue
					}
					buttonState = GetButtonSignal(def.ButtonType(b), f)
					if buttonState != lastButtonState[f][b] {
						lastButtonState[f][b] = buttonState
						elevatorEvents.Button <- def.ButtonEvent{def.Button{f, def.ButtonType(b)}, buttonState}
					}
				}
			}
		}
	}
}
