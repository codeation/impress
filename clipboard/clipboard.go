package clipboard

// Clipboard content types
const (
	UnknowType int = 0
	TextType   int = 1
)

// Clipboarder is the interface that groups clipboard contents
type Clipboarder interface {
	Type() int
	Data() []byte
}

// Text is text content of clipboard
type Text string

// Type returns clipboard content type
func (c Text) Type() int {
	return TextType
}

// Data returns clipboard raw data
func (c Text) Data() []byte {
	return []byte(c)
}

// unknown is unknown content type of clipboard
type unknown []byte

func (c unknown) Type() int    { return UnknowType }
func (c unknown) Data() []byte { return c }

// Parse returns decoded clipboard struct
func Parse(typeID int, data []byte) Clipboarder {
	switch typeID {
	case TextType:
		return Text(string(data))
	default:
		return unknown(data)
	}
}
