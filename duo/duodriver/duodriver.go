// Package to connect to GTK driver
package duodriver

import (
	"errors"
	"fmt"
	"log"

	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/joint/bus"
	"github.com/codeation/impress/joint/itrun"
	"github.com/codeation/impress/joint/lazy"
	"github.com/codeation/impress/joint/pipedriver"
)

// New returns a new GTK driver
func New() (driver.Driver, error) {
	path, err := itrun.DefaultPath()
	if errors.Is(err, itrun.ErrITNotFound) {
		log.Println("impress terminal not found")
		log.Println("set IMPRESS_TERMINAL_PATH environment variable")
		return nil, fmt.Errorf("itPath: %w", err)
	}
	pipeCreator, err := bus.NewServer(itrun.New(path))
	if err != nil {
		return nil, fmt.Errorf("bus.NewServer: %w", err)
	}
	return lazy.New(pipedriver.New(pipeCreator)), nil
}
