package def

import (
	"time"
)

const (
	NumFloors = 4
	NumTypes  = 3

	TravelTime   = 2000 * time.Millisecond
	DoorOpenTime = 3000 * time.Millisecond

	DirnDown MotorDirection = -1
	DirnStop MotorDirection = 0
	DirnUp   MotorDirection = 1

	ButtonCallUp      ButtonType = 0
	ButtonCallDown    ButtonType = 1
	ButtonCallCommand ButtonType = 2

	LightTypeUp      = 0
	LightTypeDown    = 1
	LightTypeCommand = 2
	LightTypeStop    = 3

	OrderCallUp      OrderType = 0
	OrderCallDown    OrderType = 1
	OrderCallCommand OrderType = 2

	ElevatorBehaviourIdle     = 0
	ElevatorBehaviourMoving   = 1
	ElevatorBehaviourDoorOpen = 2

	BackupFile = "backup/backups/data.json"
)
