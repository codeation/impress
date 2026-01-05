// Package implements an internal mechanism to communicate with an impress terminal.
package serversocket

import (
	"encoding/base64"
	"io"
	"time"

	"golang.org/x/net/websocket"
)

type ServerSocket struct {
	conn        *websocket.Conn
	reader      io.Reader
	waitOpening chan struct{}
	waitClosing chan struct{}
}

func New() *ServerSocket {
	s := &ServerSocket{
		waitOpening: make(chan struct{}),
		waitClosing: make(chan struct{}),
	}
	return s
}

func (s *ServerSocket) Handler() websocket.Handler {
	return websocket.Handler(s.handler)
}

func (s *ServerSocket) handler(ws *websocket.Conn) {
	s.conn = ws
	s.conn.SetDeadline(time.Now().Add(24 * time.Hour))
	s.reader = base64.NewDecoder(base64.StdEncoding, s.conn)
	close(s.waitOpening)
	<-s.waitClosing
}

func (s *ServerSocket) Read(data []byte) (int, error) {
	<-s.waitOpening
	return s.reader.Read(data)
}

func (s *ServerSocket) Write(data []byte) (int, error) {
	<-s.waitOpening
	enc := base64.NewEncoder(base64.StdEncoding, s.conn)
	length, err := enc.Write(data)
	if err != nil {
		return length, err
	}
	if err = enc.Close(); err != nil {
		return 0, err
	}
	return length, nil
}

func (s *ServerSocket) Close() error {
	defer close(s.waitClosing)
	return s.conn.Close()
}
