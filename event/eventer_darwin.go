//go:build darwin

package event

// Platform specified keyboard events
var (
	KeySave  = Keyboard{Rune: 115, Name: "s", Meta: true}
	KeyExit  = Keyboard{Rune: 113, Name: "q", Meta: true}
	KeyCopy  = Keyboard{Rune: 99, Name: "c", Meta: true}
	KeyPaste = Keyboard{Rune: 118, Name: "v", Meta: true}
)
