// Package implements an internal mechanism to communicate with an impress terminal.
package serversocket

import (
	"encoding/base64"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

type ServerSocket struct {
	conn        *websocket.Conn
	isConnected chan struct{}
	connectOnce sync.Once
	reader      io.Reader
	isClosed    chan struct{}
	closeOnce   sync.Once
}

func New() *ServerSocket {
	s := &ServerSocket{
		isConnected: make(chan struct{}),
		isClosed:    make(chan struct{}),
	}
	return s
}

func (s *ServerSocket) Handler() websocket.Handler {
	return websocket.Handler(s.handler)
}

func (s *ServerSocket) Read(data []byte) (int, error) {
	s.connectOnce.Do(s.waitConnect)
	length, err := s.reader.Read(data)
	if err != nil {
		s.closeOnce.Do(s.close)
	}
	return length, err
}

func (s *ServerSocket) Write(data []byte) (int, error) {
	s.connectOnce.Do(s.waitConnect)
	enc := base64.NewEncoder(base64.StdEncoding, s.conn)
	length, err := enc.Write(data)
	if err != nil {
		s.closeOnce.Do(s.close)
		return length, err
	}
	err = enc.Close()
	if err != nil {
		s.closeOnce.Do(s.close)
		return 0, err
	}
	return length, nil
}

func (s *ServerSocket) Close() error {
	s.connectOnce.Do(s.waitConnect)
	s.closeOnce.Do(s.close)
	return nil
}

func (s *ServerSocket) handler(ws *websocket.Conn) {
	s.conn = ws
	s.reader = base64.NewDecoder(base64.StdEncoding, s.conn)
	close(s.isConnected)
	<-s.isClosed
}

func (s *ServerSocket) waitConnect() {
	<-s.isConnected
}

func (s *ServerSocket) close() {
	close(s.isClosed)
	_ = s.conn.Close()
}
