package iface

const (
	ApplicationSizeCode    byte = 'S'
	ApplicationTitleCode   byte = 'T'
	ApplicationExitCode    byte = 'X'
	ApplicationVersionCode byte = 'V'

	FrameNewCode   byte = 'Y'
	FrameDropCode  byte = 'Q'
	FrameSizeCode  byte = 'H'
	FrameRaiseCode byte = 'J'

	WindowNewCode   byte = 'D'
	WindowDropCode  byte = 'O'
	WindowRaiseCode byte = 'A'
	WindowClearCode byte = 'C'
	WindowShowCode  byte = 'W'
	WindowSizeCode  byte = 'Z'
	WindowFillCode  byte = 'F'
	WindowLineCode  byte = 'L'
	WindowTextCode  byte = 'U'
	WindowImageCode byte = 'I'

	FontNewCode   byte = 'N'
	FontDropCode  byte = 'K'
	FontSplitCode byte = 'P'
	FontSizeCode  byte = 'R'

	ImageNewCode  byte = 'B'
	ImageDropCode byte = 'M'

	MenuNewCode  byte = 'E'
	MenuItemCode byte = 'G'

	ClipboardGetCode byte = '1'
	ClipboardPutCode byte = '2'

	EventGeneralCode   byte = 'g'
	EventKeyboardCode  byte = 'k'
	EventConfigureCode byte = 'f'
	EventButtonCode    byte = 'b'
	EventMotionCode    byte = 'm'
	EventMenuCode      byte = 'u'
	EventScrollCode    byte = 's'
	EventClipboard     byte = 'c'
)
