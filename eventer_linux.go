// +build linux

package impress

// Keyboard events
var (
	KeySave = KeyboardEvent{Rune: 115, Name: "s", Control: true}
	KeyExit = KeyboardEvent{Rune: 119, Name: "w", Control: true}
)
