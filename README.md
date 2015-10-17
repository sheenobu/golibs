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

### managed

AKA apps 2.0

managed system and component hierarchies:

	import (
		"github.com/sheenobu/golibs/managed"
		"golang.org/x/net/context"

		"time"
		"os"
		"fmt"
	)

	func main() {
		system := managed.NewSystem("MySystem")
		system.Start()

		system.RegisterForStop(os.Interrupt)

		myChannel := make(chan string)

		myTimer := managed.Timer("MyTimer", 3 * time.Second, false /*runImmediately*/, func(ctx context.Context) {
			fmt.Printf("Running process\n")
			myChannel <- "hello"
		})

		myChannelListener := managed.Simple("MyChannelListener, func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case data := <-ch:
					fmt.Printf("got signal: %s!\n", data)
				}
			}
		})

		system.Add(myTimer)
		system.Add(myChannelListener)

		system.Wait()
	}


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

