package misc

import(
	. "../def"
)

func B2I(b bool) int {
	if b {
		return 1
	}
	return 0
}

func CopyElevators(elevators Elevators) Elevators {
	newElevators := make(Elevators)
	for id,e := range elevators {
		newElevators[id] = CopyElevator(e)
	}
	return newElevators
}

func CopyElevator(elevator Elevator) Elevator {
	var newElevator Elevator
	newElevator.State = ElevatorState(elevator.State)
	newElevator.Orders = CopyOrders(elevator.Orders)
	return newElevator
}

func CopyOrders(orders Orders) Orders {
	var newOrders Orders
	for f,_ := range orders {
		for t,_ := range orders[f] {
			newOrders[f][t] = orders[f][t]
		}
	}
	return newOrders
}

func Union(differentOrders []Orders) Orders {
	var newOrders Orders
	for _,orders := range differentOrders {
		for f,_ := range orders {
			for t,_ := range orders[f] {
				newOrders[f][t] = newOrders[f][t] || orders[f][t]
			}
		}
	}
	return newOrders
}

func GlobalOrders(elevators Elevators) Orders {
	var newOrders Orders
	for _,elevator := range elevators {
		for f,floorOrders := range elevator.Orders {
			for t,order := range floorOrders {
				newOrders[f][t] = newOrders[f][t] || order
			}
		}
	}
	return newOrders
}