// Package to register GTK driver as default impress driver
package duo

import (
	"log"

	"github.com/codeation/impress"
	"github.com/codeation/impress/duo/duodriver"
)

func init() {
	driver, err := duodriver.New()
	if err != nil {
		log.Fatal(err)
	}
	impress.Register(driver)
}
