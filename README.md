# sheenobu/golibs

Code that doesn't yet deserve it's own repo (yet)

## Packages

### log

extended logging functionality via log15 and golang.org/x/net/context:

	ctx := log.NewContext(ctx, params)

	log.Log(ctx).Debug(...)
	log.Log(ctx).Error(

	l := log.Log(ctx)
	l.Debug("X")

### apps

Nested application and subprocess management


					pctx
					/
				 app
				 /
		pctx  ctx
		  \   / \
		   app  process
			 \
			 ctx
			/    \
		process  process

Usage:

	import (
		"github.com/sheenobu/golibs/apps"
		"golang.org/x/net/context"

		"time"
		"os"
		"fmt"
	)

	func main() {
		app := aps.NewApp("MyApp")
		app.Start()

		app.RegisterForStop(os.Interrupt)
		
		app.SpawnSimple("MyTimer", func(ctx context.Context) {
			for {
				select {
				case <-time.After(3 * time.Second):
					fmt.Printf("Running process\n")
				case <-ctx.Done():
					fmt.Printf("Got quit\n")
					return // Leave function
				}
			}
		})

		app.Wait() // wait for app.Stop to get called (via RegisterForStop)
	}

### dispatch

Simple pub/sub style dispatcher with named channels.

