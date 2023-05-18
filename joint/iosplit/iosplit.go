// Package implements an internal mechanism to communicate with an impress terminal.
package iosplit

import (
	"io"
	"time"
)

type IOSplit struct {
	c         chan []byte
	tail      []byte
	isData    bool
	isTimeout bool
	isEternal bool
}

func NewIOSplit() *IOSplit {
	return &IOSplit{
		c: make(chan []byte, 16),
	}
}

func (c *IOSplit) WithTimeout() *IOSplit {
	c.isTimeout = true
	return c
}

func (c *IOSplit) WithEternal() *IOSplit {
	c.isEternal = true
	return c
}

func (c *IOSplit) instantNext() error {
	select {
	case c.tail = <-c.c:
		return nil
	default:
		c.isData = false
		return io.EOF
	}
}

func (c *IOSplit) timeoutNext() error {
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()
	select {
	case c.tail = <-c.c:
		return nil
	case <-timer.C:
		c.isData = false
		return io.EOF
	}
}

func (c *IOSplit) eternalNext() error {
	var ok bool
	c.tail, ok = <-c.c
	if !ok {
		return io.EOF
	}
	return nil
}

func (c *IOSplit) Read(p []byte) (int, error) {
	if len(c.tail) == 0 {
		switch {
		case c.isData:
			if err := c.instantNext(); err != nil {
				return 0, err
			}
		case c.isTimeout:
			if err := c.timeoutNext(); err != nil {
				return 0, err
			}
		case c.isEternal:
			if err := c.eternalNext(); err != nil {
				return 0, err
			}
		default:
			if err := c.instantNext(); err != nil {
				return 0, err
			}
		}
	}

	length := copy(p, c.tail)
	c.tail = c.tail[length:]
	return length, nil
}

func (c *IOSplit) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	output := make([]byte, len(p))
	copy(output, p)
	c.c <- output
	c.isData = true
	return len(p), nil
}
