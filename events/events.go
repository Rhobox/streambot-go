package events

type Event struct {
	Type EventType
}

type EventType int

const (
	EventSubscribe = iota
	EventUnsubscribe
)
