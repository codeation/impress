// Package to connect to WebAssembly driver
package canvasdriver

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
	wg         sync.WaitGroup
}

// New returns a new WebAssembly driver
func New() (*httpDriver, error) {
	flag.Parse()

	asyncSocket := serversocket.New()
	syncSocket := serversocket.New()

	eventPipe := rpc.NewPipe(new(sync.Mutex), nil, asyncSocket)
	streamPipe := rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(asyncSocket), nil)
	syncPipe := rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(syncSocket), syncSocket)

	eventChan := eventchan.New()
	eventRecv := eventrecv.New(eventChan, eventPipe)
	client := drawsend.New(streamPipe, syncPipe)
	driver := domain.New(client, eventChan, streamPipe)
	lazyDriver := lazy.New(driver)

	fmt.Println("listening on", *listen)
	httpServer := &http.Server{
		Addr:    *listen,
		Handler: newHandlers(asyncSocket, syncSocket),
	}

	h := &httpDriver{
		Driver:     lazyDriver,
		httpServer: httpServer,
		eventRecv:  eventRecv,
	}
	h.wg.Go(func() {
		if err := h.httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println(err)
			}
		}
	})
	return h, nil
}

func (h *httpDriver) Done() {
	h.eventRecv.Done()
	h.Driver.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	h.wg.Go(func() {
		if err := h.httpServer.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	})
	h.wg.Wait()
}

type handlers struct {
	asyncSocket http.Handler
	syncSocket  http.Handler
	fileServer  http.Handler
}

func newHandlers(asyncSocket http.Handler, syncSocket http.Handler) *handlers {
	return &handlers{
		asyncSocket: asyncSocket,
		syncSocket:  syncSocket,
		fileServer:  http.FileServer(http.Dir(*dir)),
	}
}

func (s *handlers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/", "/index.html", "/main.wasm", "/wasm_exec.js":
		s.fileServer.ServeHTTP(w, req)
	case "/async":
		h := s.asyncSocket
		s.asyncSocket = http.NotFoundHandler()
		h.ServeHTTP(w, req)
	case "/sync":
		h := s.syncSocket
		s.syncSocket = http.NotFoundHandler()
		h.ServeHTTP(w, req)
	default:
		http.NotFound(w, req)
	}
}
