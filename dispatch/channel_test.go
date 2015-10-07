package dispatch

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	dispatcher := New()

	if dispatcher == nil {
		t.Errorf("Expected new dispatcher to not be nil")
	}
}

func TestSimpleNamed(t *testing.T) {
	dispatcher := New()

	if dispatcher == nil {
		t.Errorf("Expected new dispatcher to not be nil")
	}

	var wg sync.WaitGroup

	var msg string

	wg.Add(1)
	dispatcher.ConnectFn("undeliverable", func(m interface{}) {
		msg = m.(string)
		wg.Done()
	})

	ch := dispatcher.Named("hello")

	if ch == nil {
		t.Errorf("Expected ch to not be nil")
	}

	ch <- "X"
	wg.Wait()

	if msg != "X" {
		t.Errorf("Expected message to be X, was %s\n", msg)
	}
}

func TestConnect(t *testing.T) {
	var wgL sync.WaitGroup
	var wgR sync.WaitGroup

	dispatcher := New()

	dispatcher.ConnectFn("event-1", func(m interface{}) {
		wgL.Done()
	})

	dispatcher.ConnectFn("event-3", func(m interface{}) {
		wgR.Done()
	})

	dispatcher.Connect("event-2", "event-1")
	dispatcher.Connect("event-2", "event-3")

	wgL.Add(2)
	wgR.Add(2)
	dispatcher.Named("event-2") <- "hello"
	dispatcher.Named("event-2") <- "goodbye"

	wgL.Wait()
	wgR.Wait()
}

func TestDispatch(t *testing.T) {
	var wgL sync.WaitGroup

	dispatcher := New()

	dispatcher.ConnectFn("event-1", func(m interface{}) {
		wgL.Done()
	})

	wgL.Add(2)
	dispatcher.Dispatch("event-1", "")
	dispatcher.Dispatch("event-1", "")

	wgL.Wait()
}
