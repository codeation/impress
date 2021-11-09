//go:build darwin
// +build darwin

package event

// Platform specified keyboard events
var (
	KeySave = KeyboardEvent{Rune: 115, Name: "s", Meta: true}
	KeyExit = KeyboardEvent{Rune: 113, Name: "q", Meta: true}
)
