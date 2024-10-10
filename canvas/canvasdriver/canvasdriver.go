// Package to connect to WebAssembly driver
package canvasdriver

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/joint/domain"
	"github.com/codeation/impress/joint/drawsend"
	"github.com/codeation/impress/joint/eventchan"
	"github.com/codeation/impress/joint/eventrecv"
	"github.com/codeation/impress/joint/lazy"
	"github.com/codeation/impress/joint/rpc"
	"github.com/codeation/impress/joint/serversocket"
)

const defaultBufferSize = 64 * 1024

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

type httpDriver struct {
	driver.Driver
	httpServer *http.Server
	eventRecv  interface{ Done() }
}

// New returns a new WebAssembly driver
func New() *httpDriver {
	flag.Parse()

	streamSocket := serversocket.New()
	syncSocket := serversocket.New()
	eventSocket := serversocket.New()
	httpServer := newServer(streamSocket.Handler(), syncSocket.Handler(), eventSocket.Handler())

	streamBuffered := bufio.NewWriterSize(streamSocket, defaultBufferSize)

	eventPipe := rpc.NewPipe(new(sync.Mutex), nil, eventSocket)
	streamPipe := rpc.NewPipe(new(sync.Mutex), streamBuffered, nil)
	syncPipe := rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(syncSocket), syncSocket)

	eventChan := eventchan.New()
	eventRecv := eventrecv.New(eventChan, eventPipe)
	client := drawsend.New(streamPipe, syncPipe)
	driver := domain.New(client, eventChan, streamPipe)
	lazyDriver := lazy.New(driver)

	return &httpDriver{
		Driver:     lazyDriver,
		httpServer: httpServer,
		eventRecv:  eventRecv,
	}
}

func (h *httpDriver) Done() {
	h.eventRecv.Done()
	h.Driver.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	if err := h.httpServer.Shutdown(ctx); err != nil {
		log.Println(err)
	}
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
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
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
