//go:build linux
// +build linux

package event

// Platform specified keyboard events
var (
	KeySave = Keyboard{Rune: 115, Name: "s", Control: true}
	KeyExit = Keyboard{Rune: 119, Name: "w", Control: true}
)
