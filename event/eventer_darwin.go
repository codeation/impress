//go:build darwin
// +build darwin

package event

// Platform specified keyboard events
var (
	KeySave = Keyboard{Rune: 115, Name: "s", Meta: true}
	KeyExit = Keyboard{Rune: 113, Name: "q", Meta: true}
)
