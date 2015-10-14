package apps

import (
	"os"
	"os/signal"
)

// RegisterForStop registers the unix signal listener for stop
func (app *App) RegisterForStop(sig ...os.Signal) {

	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	go func() {
		<-c
		app.Stop()
	}()

}
