package apps

import (
	"io"
)

type Writer interface {
	WriteSimple(app *App, name string)
	WriteApp(app *App)

	Child() Writer
}

func TextWriter(writer io.Writer, tabs int) Writer {
	return &IoWriter{
		wr:   writer,
		tabs: tabs,
	}
}

type IoWriter struct {
	wr   io.Writer
	tabs int
}

func (wr *IoWriter) Child() Writer {
	return &IoWriter{
		wr:   wr.wr,
		tabs: wr.tabs + 1,
	}
}

func (wr *IoWriter) WriteSimple(app *App, name string) {
	for i := 0; i != wr.tabs+1; i++ {
		wr.wr.Write([]byte("\t"))
	}

	//	status := app.procs[name]

	//	if status {
	wr.wr.Write([]byte("Simple: " + name + "\n"))
	//	} else {
	//		wr.wr.Write([]byte("Simple: " + name + " [stopped]\n"))
	//	}
}

func (wr *IoWriter) WriteApp(app *App) {
	for i := 0; i != wr.tabs; i++ {
		wr.wr.Write([]byte("\t"))
	}

	wr.wr.Write([]byte("App: " + app.name + "-\n"))
}

func (app *App) WriteTree(w Writer) {
	w.WriteApp(app)

	app.lock.RLock()
	for _, ch := range app.ChildrenProcs {
		w.WriteSimple(app, ch)
	}

	child := w.Child()
	for _, ch := range app.Children {
		ch.WriteTree(child)
	}
	app.lock.RUnlock()

}
