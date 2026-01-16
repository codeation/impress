// Package implements an internal mechanism to communicate with an impress terminal.
package serversocket

import (
	"io"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type ServerSocket struct {
	conn        *websocket.Conn
	wg          sync.WaitGroup
	waitOpening chan struct{}
	pipeReader  *io.PipeReader
	pipeWriter  *io.PipeWriter
}

func New() *ServerSocket {
	pipeReader, pipeWriter := io.Pipe()
	return &ServerSocket{
		waitOpening: make(chan struct{}),
		pipeReader:  pipeReader,
		pipeWriter:  pipeWriter,
	}
}

func (s *ServerSocket) Close() error {
	s.pipeWriter.Close()
	return s.conn.Close()
}

func (s *ServerSocket) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	websocket.Handler(s.server).ServeHTTP(w, req)
}

func (s *ServerSocket) server(ws *websocket.Conn) {
	s.conn = ws
	close(s.waitOpening)
	s.wg.Go(s.readAll)
	s.wg.Wait()
}

func (s *ServerSocket) Write(data []byte) (int, error) {
	<-s.waitOpening
	err := websocket.Message.Send(s.conn, data)
	return len(data), err
}

func (s *ServerSocket) Read(data []byte) (int, error) {
	return s.pipeReader.Read(data)
}

func (s *ServerSocket) readAll() {
	for {
		var buffer []byte
		if err := websocket.Message.Receive(s.conn, &buffer); err != nil {
			s.pipeWriter.CloseWithError(err)
			return
		}
		if _, err := s.pipeWriter.Write(buffer); err != nil {
			s.pipeWriter.CloseWithError(err)
			return
		}
	}
}
