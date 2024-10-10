// Package to register GTK driver as default impress driver
package duo

import (
	"github.com/codeation/impress"
	"github.com/codeation/impress/duo/duodriver"
)

func init() {
	impress.Register(duodriver.New())
}
