package low

// #cgo pkg-config: gtk+-3.0
// #include "low.h"
import "C"

import (
	"github.com/codeation/impress"
)

var (
	ShiftKeyModifier   = C.GDK_SHIFT_MASK
	LockKeyModifier    = C.GDK_LOCK_MASK
	ControlKeyModifier = C.GDK_CONTROL_MASK
	MetaKeyModifier    = C.GDK_META_MASK
	AltKeyModifier     = C.GDK_MOD1_MASK
)

var events = make(chan impress.Eventer)
var readyChan = make(chan bool) // "Ready to GTK event" chan
var readyOk = false

func EventDequeue() impress.Eventer {
	if readyOk {
		readyChan <- true
	}
	readyOk = true
	return <-events
}

func queueEvent(event impress.Eventer) {
	guiMutex.Unlock()
	defer guiMutex.Lock()
	events <- event
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
	queueEvent(impress.DestroyEvent)
}

type KeyboardEvent struct {
	Event *C.GdkEventKey
}

func KeyRune(event *C.GdkEventKey) rune {
	return rune(C.gdk_keyval_to_unicode(event.keyval))
}

func KeyName(event *C.GdkEventKey) string {
	return C.GoString(C.gdk_keyval_name(event.keyval))
}

func KeyModifier(event *C.GdkEventKey, modifier int) bool {
	return int(event.state)&modifier != 0
}

//export KeyboardCallBack
func KeyboardCallBack(event *C.GdkEventKey) {
	queueEvent(impress.KeyboardEvent{
		Rune:    KeyRune(event),
		Name:    KeyName(event),
		Shift:   KeyModifier(event, ShiftKeyModifier),
		Control: KeyModifier(event, ControlKeyModifier),
		Alt:     KeyModifier(event, AltKeyModifier),
		Meta:    KeyModifier(event, MetaKeyModifier),
	})
}
