package slack

type Event struct {
	Type string
	Data interface{}
}

type HelloEvent struct{}

type MessageEvent struct {
	Message
}
