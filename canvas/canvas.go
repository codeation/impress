// Package to connect to WebAssembly driver
package canvas

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/codeation/impress"
	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/joint/domain"
	"github.com/codeation/impress/joint/drawsend"
	"github.com/codeation/impress/joint/eventchan"
	"github.com/codeation/impress/joint/eventrecv"
	"github.com/codeation/impress/joint/rpc"
	"github.com/codeation/impress/joint/serversocket"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

type doner interface {
	Done()
}

func init() {
	flag.Parse()

	streamSocket := serversocket.New()
	syncSocket := serversocket.New()
	eventSocket := serversocket.New()
	httpServer := newServer(streamSocket.Handler(), syncSocket.Handler(), eventSocket.Handler())

	streamBuffered := bufio.NewWriter(streamSocket)

	eventPipe := rpc.NewPipe(new(sync.Mutex), nil, eventSocket)
	streamPipe := rpc.NewPipe(new(sync.Mutex), streamBuffered, nil)
	syncPipe := rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(syncSocket), syncSocket)

	eventChan := eventchan.New()
	eventRecv := eventrecv.New(eventChan, eventPipe)
	client := drawsend.New(streamPipe, syncPipe)
	driver := domain.New(client, eventChan, streamPipe)

	impress.Register(&httpDriver{
		Driver:     driver,
		httpServer: httpServer,
		eventRecv:  eventRecv,
	})
}

type httpDriver struct {
	driver.Driver
	httpServer *http.Server
	eventRecv  doner
}

func (h *httpDriver) Done() {
	h.eventRecv.Done()
	h.Driver.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	h.httpServer.Shutdown(ctx)
}

func newServer(streamSocket, syncSocket, eventSocket websocket.Handler) *http.Server {
	s := &http.Server{
		Addr:           *listen,
		Handler:        newHandlers(streamSocket, syncSocket, eventSocket),
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("listening on", *listen)
	go s.ListenAndServe()
	return s
}

type handlers struct {
	handlers map[string]http.Handler
}

func newHandlers(streamSocket, syncSocket, eventSocket websocket.Handler) *handlers {
	fileServer := http.FileServer(http.Dir(*dir))
	s := &handlers{
		handlers: map[string]http.Handler{},
	}
	for _, path := range []string{"/", "/index.html", "/main.wasm", "/wasm_exec.js"} {
		s.handlers[path] = fileServer
	}
	s.handlers["/stream"] = streamSocket
	s.handlers["/sync"] = syncSocket
	s.handlers["/event"] = eventSocket
	return s
}

func (s *handlers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := s.handlers[req.URL.Path]; ok {
		handler.ServeHTTP(w, req)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
