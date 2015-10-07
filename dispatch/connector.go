package dispatch

// Connector provides the functionality to connect two named endpoints together
type Connector interface {
	// Connect connects two channels together
	Connect(src string, dest string)

	// ConnectFn connects a channel and a function together
	ConnectFn(src string, onMessage func(msg interface{}))
}
