package low

// #cgo pkg-config: gtk+-3.0
// #include "low.h"
import "C"

var eventChan = make(chan interface{}) // Event from GTK to golang
var readyChan = make(chan bool)        // "Ready to GTK event" chan

var readyOk = false

func (a *Application) Event() interface{} {
	if readyOk {
		readyChan <- true
	}
	readyOk = true
	return <-eventChan
}

func queueEvent(event interface{}) {
	guiMutex.Unlock()
	defer guiMutex.Lock()
	eventChan <- event
	<-readyChan
}

type GeneralEvent struct {
	Type int
}

var (
	DestroyEevent = GeneralEvent{Type: 1}
)

//export DestroyCallBack
func DestroyCallBack() {
	queueEvent(DestroyEevent)
}

type KeyboardEvent struct {
	Event *C.GdkEventKey
}

func KeyRune(e KeyboardEvent) rune {
	return rune(C.gdk_keyval_to_unicode(e.Event.keyval))
}

func KeyName(e KeyboardEvent) string {
	return C.GoString(C.gdk_keyval_name(e.Event.keyval))
}

//export KeyboardCallBack
func KeyboardCallBack(event *C.GdkEventKey) {
	queueEvent(KeyboardEvent{Event: event})
}
