package dispatch

// Dispatcher provides named publishing of messages
type Dispatcher interface {
	// Dispatch sends a message to a named channel
	Dispatch(name string, message interface{})

	// TryDispatch attempts to send a message to a named channel, failing if the message is undeliverable
	TryDispatch(name string, message interface{}) error
}
