package sink

type Sink interface {
	Send(payload *Payload) error
}
