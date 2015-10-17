package managed

import (
	"io"
)

// Writer defines the interface for writing system data
type Writer interface {

	// WriteProcess writes the process
	WriteProcess(sys *System, process *Process)

	// WriteSystem writes the system
	WriteSystem(sys *System)

	// Child gets the child writer
	Child() Writer
}

// TextWriter returns the plaintext writer
func TextWriter(writer io.Writer, tabs int) Writer {
	return &IoWriter{
		wr:   writer,
		tabs: tabs,
	}
}

// IoWriter is the plaintext writer that write to io.Writer
type IoWriter struct {
	wr   io.Writer
	tabs int
}

// Child returns the child IoWriter with the tab spacing shifted over by 1
func (wr *IoWriter) Child() Writer {
	return &IoWriter{
		wr:   wr.wr,
		tabs: wr.tabs + 1,
	}
}

// WriteProcess writes the process information
func (wr *IoWriter) WriteProcess(sys *System, process *Process) {
	for i := 0; i != wr.tabs+1; i++ {
		wr.wr.Write([]byte("\t"))
	}

	process.Writer(process, wr.wr)
}

// WriteSystem writes the system information
func (wr *IoWriter) WriteSystem(sys *System) {
	for i := 0; i != wr.tabs; i++ {
		wr.wr.Write([]byte("\t"))
	}

	wr.wr.Write([]byte("System: " + sys.name + "-\n"))
}

// WriteTree recursively writes the system(s) and processes to the given writer.
func (sys *System) WriteTree(w Writer) {
	w.WriteSystem(sys)

	sys.lock.RLock()
	for _, ch := range sys.ChildrenProcs {
		w.WriteProcess(sys, ch)
	}

	child := w.Child()
	for _, ch := range sys.Children {
		ch.WriteTree(child)
	}
	sys.lock.RUnlock()

}
