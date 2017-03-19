package network


type Buffer struct {
	StateMessages []StateMessage
	OrderMessages []OrderMessage
}

func NewBuffer() Buffer {
    return Buffer{make([]StateMessage, 0), make([]OrderMessage, 0)}
}

func (b *Buffer) EnqueueStateMessage(stateMessage StateMessage) {
	b.StateMessages = append(b.StateMessages, stateMessage)
}

func (b *Buffer) EnqueueOrderMessage(orderMessage OrderMessage) {
	b.OrderMessages = append(b.OrderMessages, orderMessage)
}

func (b *Buffer) DequeueStateMessage() {
	b.StateMessages = b.StateMessages[1:]
}

func (b *Buffer) DequeueOrderMessage() {
	b.OrderMessages = b.OrderMessages[1:]
}

func (b *Buffer) TopStateMessage() StateMessage {
	return b.StateMessages[0]
}

func (b *Buffer) TopOrderMessage() OrderMessage {
	return b.OrderMessages[0]
}

func (b *Buffer) HasStateMessage() bool {
	return len(b.StateMessages) > 0
}

func (b *Buffer) HasOrderMessage() bool {
	return len(b.OrderMessages) > 0
}