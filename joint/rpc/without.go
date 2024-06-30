package rpc

import (
	"io"
)

type withoutMutex struct{}

func WithoutMutex() *withoutMutex {
	return &withoutMutex{}
}

func (*withoutMutex) Lock()   {}
func (*withoutMutex) Unlock() {}

type withoutFlush struct {
	io.Writer
}

func WithoutFlush(w io.Writer) *withoutFlush {
	return &withoutFlush{
		Writer: w,
	}
}

func (w *withoutFlush) Flush() error { return nil }
