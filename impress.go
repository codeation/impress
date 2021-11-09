package impress

import (
	"github.com/codeation/impress/driver"
)

var d driver.Driver

// Register is an internal function that makes the GUI driver available.
func Register(current driver.Driver) {
	d = current
}
