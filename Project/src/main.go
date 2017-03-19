package main

import (
	. "./def"
	"./driver"
	"./elevator"
	//"./gui"
	"./network"
	//"./network/conn"
	"./gui"
	"./network/localip"
	"./orders"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	// Handle terminal input: id and simulation config file
	var id string
	var simulator string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.StringVar(&simulator, "sim", "", "simulator config file")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	if simulator == "" {
		driver.Init(driver.TypeComedi, simulator)
	} else {
		driver.Init(driver.TypeSimulation, simulator)
	}

	// Initialize channels and routines
	driverElevatorEvents := DriverElevatorEvents{
		Button:         make(chan ButtonEvent, 10),
		Light:          make(chan LightEvent, 10),
		Stop:           make(chan bool, 10),
		MotorDirection: make(chan MotorDirection, 10),
		Floor:          make(chan int, 10),
		DoorOpen:       make(chan bool, 10),
		FloorIndicator: make(chan int, 10),
	}

	elevatorOrdersEvents := ElevatorOrdersEvents{
		Order:        make(chan Order, 10),
		State:        make(chan ElevatorState, 10),
		LocalOrders:  make(chan Orders, 10),
		GlobalOrders: make(chan Orders, 10),
	}

	ordersNetworkEvents := OrdersNetworkEvents{
		TxOrderEvent: make(chan OrderEvent, 10),
		RxOrderEvent: make(chan OrderEvent, 10),
		TxStateEvent: make(chan StateEvent, 10),
		RxStateEvent: make(chan StateEvent, 10),
		ElevatorNew:  make(chan string, 10),
		ElevatorLost: make(chan string, 10),
		Elevators:    make(chan Elevators, 10),
	}

	ordersGuiEvents := OrdersGuiEvents{
		Elevators: make(chan Elevators, 10),
	}

	go network.EventManager(id, ordersNetworkEvents)
	go elevator.EventManager(driverElevatorEvents, elevatorOrdersEvents)
	go driver.EventManager(driverElevatorEvents)
	go orders.OrderManager(id, elevatorOrdersEvents, ordersNetworkEvents, ordersGuiEvents)
	go gui.ElevatorVisualizer(id, ordersGuiEvents)

	for {
		select {
		case <-time.After(5 * time.Second):
		}
	}
}
