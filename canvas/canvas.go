// Package to register WebAssembly driver as default impress driver
package canvas

import (
	"log"

	"github.com/codeation/impress"
	"github.com/codeation/impress/canvas/canvasdriver"
)

func init() {
	driver, err := canvasdriver.New()
	if err != nil {
		log.Fatal(err)
	}
	impress.Register(driver)
}
