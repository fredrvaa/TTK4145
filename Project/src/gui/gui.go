package gui

import (
	. "../def"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"time"
)

func clearTerminal() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func sortElevators(elevs Elevators) []string {
	mk := make([]string, len(elevs))
	i := 0
	for k, _ := range elevs {
		mk[i] = k
		i++
	}
	sort.Strings(mk)
	return mk
}

func print(id string, elevs Elevators) {
	for _, elev := range sortElevators(elevs) {
		if elev == id {
			fmt.Println("ELEVATOR", elev, "(this)")
		} else {
			fmt.Println("ELEVATOR", elev)
		}
		fmt.Println("Floor	State		Up	Down	Command")
		for f := len(elevs[elev].Orders) - 1; f >= 0; f-- {
			fmt.Print(f, "	")

			if f == len(elevs[elev].Orders)-1 && elevs[elev].State.Floor < 0 {
				fmt.Print("U	")
			} else if f == elevs[elev].State.Floor {
				if elevs[elev].State.Active {
					fmt.Print("(A) ")
				} else {
					fmt.Print("(I) ")
				}
				switch elevs[elev].State.Direction {
				case DirnUp:
					fmt.Print("ˆ ")
				case DirnDown:
					fmt.Print("v ")
				case DirnStop:
					fmt.Print("  ")
				}
				switch elevs[elev].State.Behaviour {
				case ElevatorBehaviourIdle:
					fmt.Print(" []")
				case ElevatorBehaviourMoving:
					fmt.Print(" []*")
				case ElevatorBehaviourDoorOpen:
					fmt.Print("[  ]")
				}
			} else {
				fmt.Print("	")
			}
			for t, order := range elevs[elev].Orders[f] {
				fmt.Print("	")
				if f == len(elevs[elev].Orders)-1 && t == int(OrderCallUp) {
					continue
				}
				if f == 0 && t == int(OrderCallDown) {
					continue
				}
				if order {
					fmt.Print("*")
				} else {
					fmt.Print("-")
				}
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func ElevatorVisualizer(id string, ordersEvents OrdersGuiEvents) {
	elevs := make(Elevators)
	clearTerminal()
	print(id, elevs)

	for {
		select {
		case elevs = <-ordersEvents.Elevators:
			clearTerminal()
			print(id, elevs)
		case <-time.After(100 * time.Millisecond):
		}
	}
}
