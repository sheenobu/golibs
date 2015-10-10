package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sheenobu/golibs/apps"
	"golang.org/x/net/context"
)

func main() {

	app := apps.NewApp("app1")
	app.Start()

	sender := make(chan int)
	accum := make(chan int)

	go func() {
		a := 0
		for i := 0; i < 5; i++ {
			select {
			case val := <-accum:
				a += val
			}
		}

		fmt.Printf("Got: %d\n", a)
		app.Stop()
	}()

	for _, i := range []int{1, 2, 3, 4, 5} {
		name := "my." + strconv.FormatInt(int64(i), 10)
		app.SpawnSimple("my."+strconv.FormatInt(int64(i), 10), mySubprocessFactory(name, sender, accum))
	}

	go func() {
		<-time.After(1 * time.Millisecond)
		app.WriteTree(apps.TextWriter(os.Stdout, 1))
	}()

	go func() {
		sender <- 1
		sender <- 2
		sender <- 3
		sender <- 4
		sender <- 5
	}()

	app.Wait()

	<-time.After(3 * time.Millisecond)
}

func mySubprocessFactory(name string, recv <-chan int, resp chan<- int) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("worker %s quitting\n", name)
				return
			case i := <-recv:
				fmt.Printf("worker %s answering\n", name)
				resp <- i * i
			}
		}
	}
}
