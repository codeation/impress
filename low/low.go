package low

// #cgo pkg-config: gtk+-3.0
// #include "low.h"
import "C"

import (
	"image"
	"sync"
)

var guiMutex sync.Mutex

type Application struct {
	window *C.GtkWidget
	layout *C.GtkWidget
}

func NewApplication() *Application {
	guiMutex.Lock()
	defer guiMutex.Unlock()
	app := C.application_create()
	layout := C.layout_create(app)
	return &Application{
		window: app,
		layout: layout,
	}
}

func (a *Application) Main() {
	guiMutex.Lock()
	defer guiMutex.Unlock()
	C.application_main(a.window)
}

func (a *Application) Quit() {
	guiMutex.Lock()
	defer guiMutex.Unlock()
	C.application_quit()
	readyChan <- true
}

func (a *Application) Size(x, y, width, height int) {
	guiMutex.Lock()
	defer guiMutex.Unlock()
	C.application_size(a.window, C.int(x), C.int(y), C.int(width), C.int(height))
}

type Window struct {
	window *C.GtkWidget
}

func (a *Application) NewWindow() *Window {
	guiMutex.Lock()
	defer guiMutex.Unlock()
	return &Window{
		window: C.window_create(a.layout),
	}
}

func (w *Window) Move(a *Application, x, y int) {
	guiMutex.Lock()
	defer guiMutex.Unlock()
	C.window_move(a.layout, w.window, C.int(x), C.int(y))
}

func (w *Window) Set(img *image.RGBA) {
	guiMutex.Lock()
	defer guiMutex.Unlock()
	C.window_set(w.window, C.int(img.Rect.Size().X), C.int(img.Rect.Size().Y),
		C.int(img.Stride), C.CBytes(img.Pix))
}
