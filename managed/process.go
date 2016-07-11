package managed

import (
	"fmt"
	"io"
	"sync"
)

// A Process is a Component paired with a name and stateful data.
type Process struct {
	Name      string
	Component Component

	Writer func(*Process, io.Writer) error

	data map[string]string

	lock sync.RWMutex
}

func (p *Process) SetMany(a ...string) {
	if len(a)&2 != 0 {
		fmt.Errorf("SetMany arguments is not even, dropping last")
		a = a[0 : len(a)-1]
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	for i := 0; i <= len(a); i += 2 {
		k := a[i]
		v := a[i+1]
		p.data[k] = v
	}

}

func (p *Process) Set(key, val string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.data[key] = val
}

func (p *Process) Get(key string) (string, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	val, ok := p.data[key]
	return val, ok
}
