//go:build linux
// +build linux

package event

// Platform specified keyboard events
var (
	KeySave  = Keyboard{Rune: 115, Name: "s", Control: true}
	KeyExit  = Keyboard{Rune: 119, Name: "w", Control: true}
	KeyCopy  = Keyboard{Rune: 99, Name: "c", Control: true}
	KeyPaste = Keyboard{Rune: 118, Name: "v", Control: true}
)
