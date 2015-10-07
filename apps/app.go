package apps

import (
	"golang.org/x/net/context"

	log "gopkg.in/inconshreveable/log15.v2"

	"sync"
	"time"
)

// App is our toplevel application construct
type App struct {
	name string

	parentContext context.Context
	parentCancel  func()

	ctx      context.Context
	cancelFn func()
	wg       sync.WaitGroup

	startChan chan string
	stopChan  chan string
}

// NewApp creates a new application
func NewApp(name string) *App {
	return &App{
		name:      name,
		startChan: make(chan string),
		stopChan:  make(chan string),
	}
}

// Start starts the application
func (app *App) Start() {

	log.Debug("Starting application", "app", app.name)

	app.parentContext = context.Background()
	app.parentContext, app.parentCancel = context.WithCancel(app.parentContext)

	app.ctx, app.cancelFn = context.WithCancel(app.parentContext)

}

// StartWithParent starts the application with the parent app as the context
func (app *App) StartWithParent(parent *App) {

	log.Debug("Starting application with parent application", "app", app.name)

	app.parentContext = context.Background()
	app.parentContext, app.parentCancel = context.WithCancel(app.parentContext)

	app.ctx, app.cancelFn = context.WithCancel(parent.ctx)
}

// Stop stops the application
func (app *App) Stop() {

	log.Debug("Stopping application", "app", app.name)

	app.cancelFn()
}

// SpawnSimple starts a simple single argument function as a managed process of this application
func (app *App) SpawnSimple(name string, f func(ctx context.Context)) {
	go func() {
		app.startChan <- name
		f(app.ctx)
		app.stopChan <- name
	}()
}

// SpawnApp starts an application as a managed process of this application
func (app *App) SpawnApp(child *App) {
	go func() {
		app.startChan <- "app:" + child.name
		child.StartWithParent(app)
		child.Wait()
		app.stopChan <- "app:" + child.name
	}()
}

// Wait waits for the application and its subprocesses to stop
func (app *App) Wait() {
	log.Debug("Waiting on application stop", "app", app.name)

	procs := make(map[string]bool)

	go func() {
		for {
			select {
			case proc := <-app.startChan:
				log.Debug("Got suprocess start", "process", proc, "app", app.name)
				app.wg.Add(1)
				procs[proc] = true
			case proc := <-app.stopChan:
				log.Debug("Got subprocess stop", "process", proc, "app", app.name)
				app.wg.Done()
				procs[proc] = false
			case <-app.ctx.Done():
				return
			}
		}
	}()

	<-app.ctx.Done()

	go func() {
		for {
			select {
			case proc := <-app.stopChan:
				log.Debug("Got subprocess stop", "process", proc, "app", app.name)
				procs[proc] = false
				app.wg.Done()
			case <-app.parentContext.Done():
				return
			}
		}
	}()

	ch := make(chan struct{})

	log.Debug("Application stopped, waiting on subprocesses", "app", app.name)
	go func() {
		app.wg.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		log.Debug("All subprocesses stopped", "app", app.name)
	case <-time.After(1 * time.Second):
		log.Error("Some subprocesses failed to stop in time", "app", app.name)
		for k, v := range procs {
			if v == true {
				log.Error("Subprocess failed to stop in time", "proc", k, "app", app.name)
			}
		}
	}

	app.parentCancel()

}
