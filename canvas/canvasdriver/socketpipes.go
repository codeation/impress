package canvasdriver

import (
	"bufio"
	"log"
	"net/http"

	"github.com/codeation/impress/joint/rpc"
	"github.com/codeation/impress/joint/serversocket"
)

type socketPipes struct {
	asyncSocket *serversocket.ServerSocket
	syncSocket  *serversocket.ServerSocket
	fileServer  http.Handler
}

func newSocketPipes(fileServer http.Handler) *socketPipes {
	return &socketPipes{
		asyncSocket: serversocket.New(),
		syncSocket:  serversocket.New(),
		fileServer:  fileServer,
	}
}

func (p *socketPipes) done() {
	if err := p.asyncSocket.Close(); err != nil {
		log.Printf("asyncSocket.Close: %v", err)
	}
	if err := p.syncSocket.Close(); err != nil {
		log.Printf("syncSocket.Close: %v", err)
	}
}

func (p *socketPipes) newEventPipe() *rpc.Pipe {
	return rpc.NewPipe(nil, p.asyncSocket)
}

func (p *socketPipes) newStreamPipe() *rpc.Pipe {
	return rpc.NewPipe(bufio.NewWriter(p.asyncSocket), nil)
}

func (p *socketPipes) newSyncPipe() *rpc.Pipe {
	return rpc.NewPipe(bufio.NewWriter(p.syncSocket), p.syncSocket)
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
