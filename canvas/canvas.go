// Package to register WebAssembly driver as default impress driver
package canvas

import (
	"github.com/codeation/impress"
	"github.com/codeation/impress/canvas/canvasdriver"
)

func init() {
	impress.Register(canvasdriver.New())
}
