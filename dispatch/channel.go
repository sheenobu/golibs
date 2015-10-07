package dispatch

// ChannelSystem is a channel based dispatching system.
type ChannelSystem interface {
	Dispatcher
	Connector

	// Named creates and returns a named channel, if it doesn't exist
	Named(name string) chan<- interface{}
}

// New creates a new ChannelSystem
func New() ChannelSystem {
	dispatcher := &simpleDispatcherType{
		data: make(map[string]*namedChannel),
	}
	return dispatcher
}
