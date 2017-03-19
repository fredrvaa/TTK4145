package def

// EventManager
type Button struct {
	Floor int
	Type ButtonType
}

type ButtonEvent struct {
	Button Button
	State  bool
}

type LightType int

type LightEvent struct {
	LightType LightType
	Floor     int
	Value     bool
}

// Orders
type Order struct {
	Floor int
	Type  OrderType
	Flag  bool
}

type Orders [NumFloors][NumTypes]bool


type OrderType int

// Elevator
type ElevatorBehaviour int

type ElevatorState struct {
	Active	  bool
	Floor     int
	Direction MotorDirection
	Behaviour ElevatorBehaviour
}

type Elevator struct {
	State  ElevatorState
	Orders Orders
}

type Elevators map[string]Elevator

// Driver
type MotorDirection int

type ButtonType int

type DriverElevatorEvents struct {
	Button chan ButtonEvent
	Light chan LightEvent
	Stop chan bool
	MotorDirection chan MotorDirection
	Floor chan int
	DoorOpen chan bool
	FloorIndicator chan int
}

//Network
type OrderEvent struct {
	Target string
	Order Order
}

type StateEvent struct {
	Target string
	State ElevatorState
}

type ElevatorOrdersEvents struct {
	Order chan Order
	State chan ElevatorState
	LocalOrders chan Orders
	GlobalOrders chan Orders
}

type OrdersNetworkEvents struct {
	TxOrderEvent chan OrderEvent
	RxOrderEvent chan OrderEvent
	TxStateEvent chan StateEvent
	RxStateEvent chan StateEvent
	ElevatorNew chan string
	ElevatorLost chan string
	Elevators chan Elevators
}

type MessageElevator struct {
	Id string
	Elevator Elevator
}

//GUI
type OrdersGuiEvents struct {
	Elevators chan Elevators
}



