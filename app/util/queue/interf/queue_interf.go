package interf

type QueueInterface interface {
	Handle(interface{})
}

type QueuePusherInterface interface {
	Send(string, int) bool
}
