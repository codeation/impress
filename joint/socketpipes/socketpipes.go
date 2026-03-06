// Package implements an internal mechanism to communicate with an impress terminal.
package socketpipes

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/codeation/impress/joint/rpc"
)

type Config struct {
	Listen      string
	FS          fs.FS
	NewSocketFn func() SocketHandler
}

type SocketHandler interface {
	io.ReadWriteCloser
	http.Handler
}

type socketPipes struct {
	asyncSocket SocketHandler
	syncSocket  SocketHandler
	fileServer  http.Handler
	httpServer  *http.Server
	wg          sync.WaitGroup
}

func NewSocketPipes(cfg *Config) *socketPipes {
	p := &socketPipes{
		asyncSocket: cfg.NewSocketFn(),
		syncSocket:  cfg.NewSocketFn(),
		fileServer:  http.FileServerFS(cfg.FS),
	}
	p.httpServer = &http.Server{
		Addr:    cfg.Listen,
		Handler: p,
	}
	p.wg.Go(func() {
		fmt.Println("listening on", p.httpServer.Addr)
		if err := p.httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Printf("httpServer.ListenAndServe: %v", err)
			}
		}
	})
	return p
}

func (p *socketPipes) NewEventPipe() *rpc.Pipe {
	return rpc.NewPipe(nil, p.asyncSocket)
}

func (p *socketPipes) NewStreamPipe() *rpc.Pipe {
	return rpc.NewPipe(bufio.NewWriter(p.asyncSocket), nil)
}

func (p *socketPipes) NewSyncPipe() *rpc.Pipe {
	return rpc.NewPipe(bufio.NewWriter(p.syncSocket), p.syncSocket)
}

func (p *socketPipes) Close() error {
	if err := p.asyncSocket.Close(); err != nil {
		return fmt.Errorf("asyncSocket.Close: %w", err)
	}
	if err := p.syncSocket.Close(); err != nil {
		return fmt.Errorf("syncSocket.Close: %w", err)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	p.wg.Go(func() {
		if err := p.httpServer.Shutdown(ctx); err != nil {
			log.Printf("httpServer.Shutdown: %v", err)
		}
	})
	p.wg.Wait()
	return nil
}

func (p *socketPipes) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/", "/index.html", "/main.wasm", "/wasm_exec.js":
		p.fileServer.ServeHTTP(w, req)
	case "/async":
		p.asyncSocket.ServeHTTP(w, req)
	case "/sync":
		p.syncSocket.ServeHTTP(w, req)
	default:
		http.NotFound(w, req)
	}
}
