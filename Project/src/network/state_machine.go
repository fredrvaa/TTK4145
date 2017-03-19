package network

import (
	. "../def"
)

type StateMachine struct {
	id string 
	OrdersEvents OrdersNetworkEvents
	TxStateMessageCh chan StateMessage
	TxOrderMessageCh chan OrderMessage
	TxStateAckCh chan StateAck
	TxOrderAckCh chan OrderAck

	Buffer Buffer
	StateAcks map[string]bool
	OrderAcks map[string]bool
	MessageId int
	Elevators Elevators
}


func NewStateMachine(id string, ordersEvents OrdersNetworkEvents, txStateMessageCh chan StateMessage, txOrderMessageCh chan OrderMessage, txStateAckCh chan StateAck, txOrderAckCh chan OrderAck) *StateMachine {
	sm := new(StateMachine)
	sm.id = id
	sm.OrdersEvents = ordersEvents
	sm.TxStateMessageCh = txStateMessageCh
	sm.TxOrderMessageCh = txOrderMessageCh
	sm.TxStateAckCh = txStateAckCh
	sm.TxOrderAckCh = txOrderAckCh

	sm.Buffer = NewBuffer()
	sm.StateAcks = make(map[string]bool)
	sm.OrderAcks = make(map[string]bool)
	sm.MessageId = 0
	sm.Elevators = make(Elevators)
	return sm
}

func (this *StateMachine) OnStateEventTransmit(stateEvent StateEvent) {
	this.Buffer.EnqueueStateMessage(StateMessage{this.id, this.id+string(this.MessageId), stateEvent})
	this.MessageId++
}

func (this *StateMachine) OnOrderEventTransmit(orderEvent OrderEvent) {
	this.Buffer.EnqueueOrderMessage(OrderMessage{this.id, this.id+string(this.MessageId), orderEvent})
	this.MessageId++
}

func (this *StateMachine) OnStateMessageReceived(stateMessage StateMessage) {
	this.OrdersEvents.RxStateEvent <-stateMessage.StateEvent
	this.TxStateAckCh <-StateAck{this.id, stateMessage.Id}
}

func (this *StateMachine) OnOrderMessageReceived(orderMessage OrderMessage) {
	this.OrdersEvents.RxOrderEvent <-orderMessage.OrderEvent
	this.TxOrderAckCh <-OrderAck{this.id, orderMessage.Id}
}

func (this *StateMachine) OnStateAckReceived(stateAck StateAck) {
	if this.Buffer.HasStateMessage() {
		if stateAck.Id == this.Buffer.TopStateMessage().Id {
			this.StateAcks[stateAck.Source] = true
		}
	}
}

func (this *StateMachine) OnOrderAckReceived(orderAck OrderAck) {
	if this.Buffer.HasOrderMessage() {
		if orderAck.Id == this.Buffer.TopOrderMessage().Id {
			this.OrderAcks[orderAck.Source] = true
		}
	}
}

func (this *StateMachine) OnPeerNew(peer string) {
	this.OrdersEvents.ElevatorNew <-peer
	this.Buffer.EnqueueStateMessage(StateMessage{this.id, this.id+string(this.MessageId), StateEvent{this.id, this.Elevators[this.id].State}})
	this.MessageId++
	for f,_ := range this.Elevators[this.id].Orders {
		for t,_ := range this.Elevators[this.id].Orders[f] {
			if this.Elevators[this.id].Orders[f][t] {
				this.Buffer.EnqueueOrderMessage(OrderMessage{this.id, this.id+string(this.MessageId), OrderEvent{this.id, Order{f,OrderType(t), true}}})
				this.MessageId++
			}
		}
	}
}

func (this *StateMachine) OnPeerLost(peer string) {
	this.OrdersEvents.ElevatorLost <-peer
}

func (this *StateMachine) OnElevatorsUpdated(elevators Elevators) {
	this.Elevators = elevators
}

func (this *StateMachine) OnInterval() {
	if this.Buffer.HasStateMessage() {
		dequeueState := true
		for e,_ := range this.Elevators {
			if this.id == e {
				continue
			}
			_, ok := this.StateAcks[e]
	        if !ok {
	        	dequeueState = false
	        	break
	        }
		}
		if dequeueState {
			this.Buffer.DequeueStateMessage()
			this.StateAcks = make(map[string]bool)
		} else {
			for i := 0; i < 3; i++ {
				this.TxStateMessageCh <- this.Buffer.TopStateMessage()
			}
		}
	}
	
	if this.Buffer.HasOrderMessage() {
		dequeueOrder := true
		for e,_ := range this.Elevators {
			if this.id == e {
				continue
			}
			_, ok := this.OrderAcks[e]
	        if !ok {
	        	dequeueOrder = false
	        	break
	        }
		}
		if dequeueOrder {
			this.Buffer.DequeueOrderMessage()
			this.OrderAcks = make(map[string]bool)
		} else {
			for i := 0; i < 3; i++ {
				this.TxOrderMessageCh <- this.Buffer.TopOrderMessage()
			}
		}
	}
}