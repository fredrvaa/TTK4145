package network

import (
	. "../def"
	"./peers"
	"./bcast"
	"time"
)

const interval = 15 * time.Millisecond

type OrderMessage struct {
	Source string
	Id string
	OrderEvent OrderEvent
}
type StateMessage struct {
	Source string
	Id string
	StateEvent StateEvent
}

type StateAck struct {
	Source string
	Id string
}
type OrderAck struct {
	Source string
	Id string
}

func EventManager(id string, ordersEvents OrdersNetworkEvents) {
	
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	txStateMessageCh := make(chan StateMessage)
	rxStateMessageCh := make(chan StateMessage)

	txOrderMessageCh := make(chan OrderMessage)
	rxOrderMessageCh := make(chan OrderMessage)

	txStateAckCh := make(chan StateAck)
	rxStateAckCh := make(chan StateAck)

	txOrderAckCh := make(chan OrderAck)
	rxOrderAckCh := make(chan OrderAck)

	go peers.Transmitter(20004, id, peerTxEnable)
	go peers.Receiver(20004, peerUpdateCh)

	go bcast.Transmitter(26004, txStateMessageCh, txOrderMessageCh, txStateAckCh, txOrderAckCh)
	go bcast.Receiver(26004, rxStateMessageCh, rxOrderMessageCh, rxStateAckCh, rxOrderAckCh)

	sm := NewStateMachine(id, ordersEvents, txStateMessageCh, txOrderMessageCh, txStateAckCh, txOrderAckCh)

	for {
		select {

		case stateEvent := <-ordersEvents.TxStateEvent:
			sm.OnStateEventTransmit(stateEvent)

		case orderEvent := <-ordersEvents.TxOrderEvent:
			sm.OnOrderEventTransmit(orderEvent)

		case stateMessage := <-rxStateMessageCh:
			if stateMessage.Source != id {
				sm.OnStateMessageReceived(stateMessage)
			}

		case orderMessage := <-rxOrderMessageCh:
			if orderMessage.Source != id {
				sm.OnOrderMessageReceived(orderMessage)
			}

		case stateAck := <-rxStateAckCh:
			if stateAck.Source != id {
				sm.OnStateAckReceived(stateAck)
			}

		case orderAck := <-rxOrderAckCh:
			if orderAck.Source != id {
				sm.OnOrderAckReceived(orderAck)
			}

		case peerUpdate := <-peerUpdateCh:
			if peerUpdate.New != "" {
				if peerUpdate.New != id {
					sm.OnPeerNew(peerUpdate.New)
				}
			}
			for _,lostElevator := range peerUpdate.Lost {
				if lostElevator != id {
					sm.OnPeerLost(lostElevator)
				}
			}

		case elevators := <-ordersEvents.Elevators:
			sm.OnElevatorsUpdated(elevators)

		case <-time.After(interval):
			sm.OnInterval()
		}
	}
}