// Package to connect to WebAssembly driver
package canvasdriver

import (
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
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

type httpDriver struct {
	driver.Driver
	pipes      *socketPipes
	httpServer *http.Server
	eventRecv  interface{ Done() }
	wg         sync.WaitGroup
}

// New returns a new WebAssembly driver
func New() (*httpDriver, error) {
	flag.Parse()

	pipes := newSocketPipes(http.FileServer(http.Dir(*dir)))

	eventPipe := pipes.newEventPipe()
	streamPipe := pipes.newStreamPipe()
	syncPipe := pipes.newSyncPipe()

	eventChan := eventchan.New()
	eventRecv := eventrecv.New(eventChan, eventPipe)
	client := drawsend.New(streamPipe, syncPipe)
	driver := domain.New(client, eventChan, streamPipe)
	lazyDriver := lazy.New(driver)

	h := &httpDriver{
		Driver: lazyDriver,
		pipes:  pipes,
		httpServer: &http.Server{
			Addr:    *listen,
			Handler: pipes,
		},
		eventRecv: eventRecv,
	}
	h.wg.Go(func() {
		fmt.Println("listening on", h.httpServer.Addr)
		if err := h.httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Printf("httpServer.ListenAndServe: %v", err)
			}
		}
	})
	return h, nil
}

func (h *httpDriver) Done() {
	h.eventRecv.Done()
	h.pipes.done()
	h.Driver.Done()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	h.wg.Go(func() {
		if err := h.httpServer.Shutdown(ctx); err != nil {
			log.Printf("httpServer.Shutdown: %v", err)
		}
	})
	h.wg.Wait()
}
