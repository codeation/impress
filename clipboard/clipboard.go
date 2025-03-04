package clipboard

// Clipboard content types.
const (
	UnknowType int = 0 // UnknowType represents an unknown clipboard content type
	TextType   int = 1 // TextType represents a text clipboard content type
)

// Clipboarder is the interface that groups clipboard contents.
type Clipboarder interface {
	Type() int    // Type returns the type of clipboard content
	Data() []byte // Data returns the raw data of the clipboard content
}

// Text represents the text content of the clipboard.
type Text string

// Type returns the clipboard content type for text.
func (c Text) Type() int {
	return TextType
}

// Data returns the raw data of the text clipboard content.
func (c Text) Data() []byte {
	return []byte(c)
}

// unknown represents an unknown content type of the clipboard.
type unknown []byte

// Type returns the clipboard content type for unknown content.
func (c unknown) Type() int {
	return UnknowType
}

// Data returns the raw data of the unknown clipboard content.
func (c unknown) Data() []byte {
	return c
}

// Parse decodes the clipboard content based on the typeID and returns the appropriate Clipboarder.
func Parse(typeID int, data []byte) Clipboarder {
	switch typeID {
	case TextType:
		return Text(string(data))
	default:
		return unknown(data)
	}
}
