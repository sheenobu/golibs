package dispatch

import (
	"fmt"
)

type chanList []chan<- interface{}

type namedChannel struct {
	name     string
	incoming chan interface{}
	channels chanList
}

func (nc *namedChannel) execute(undeliverable chan interface{}) {
	go func() {
		for msg := range nc.incoming {
			if len(nc.channels) == 0 && undeliverable != nil {
				go func(msg interface{}) {
					undeliverable <- msg
				}(msg)
				continue
			}
			for _, ch := range nc.channels {
				go func(ch chan<- interface{}, msg interface{}) {
					ch <- msg
				}(ch, msg)
			}
		}
	}()
}

type simpleDispatcherType struct {
	data map[string]*namedChannel
}

func (disp *simpleDispatcherType) Named(name string) chan<- interface{} {
	return disp.named(name).incoming
}

func (disp *simpleDispatcherType) named(name string) *namedChannel {
	l, ok := disp.data[name]
	if !ok {
		l = &namedChannel{
			name:     name,
			incoming: make(chan interface{}),
			channels: make([]chan<- interface{}, 0),
		}
		disp.data[name] = l
		if name != "undeliverable" {
			l.execute(disp.named("undeliverable").incoming)
		} else {
			l.execute(nil)
		}
	}
	return l
}

func (disp *simpleDispatcherType) exists(name string) bool {
	_, ok := disp.data[name]
	return ok
}

func (disp *simpleDispatcherType) Connect(src string, dest string) {
	srcC := disp.named(src)
	dstChannel := disp.Named(dest)
	srcC.channels = append(srcC.channels, dstChannel)
}

func (disp *simpleDispatcherType) ConnectFn(src string, onMessage func(msg interface{})) {
	srcC := disp.named(src)
	dstChannel := make(chan interface{})
	srcC.channels = append(srcC.channels, dstChannel)
	go func() {
		for msg := range dstChannel {
			onMessage(msg)
		}
	}()
}

func (disp *simpleDispatcherType) Dispatch(name string, msg interface{}) {
	l := disp.named(name)
	l.incoming <- msg
}

func (disp *simpleDispatcherType) TryDispatch(name string, msg interface{}) error {
	if !disp.exists(name) {
		return fmt.Errorf("Named channel not found")
	}

	l := disp.named(name)
	if len(l.channels) == 0 {
		return fmt.Errorf("No Destinations for named channel")
	}

	disp.Dispatch(name, msg)
	return nil
}
