// Package to connect to WebAssembly driver
package canvasdriver

import (
	"flag"
	"os"

	"github.com/codeation/impress/driver"
	"github.com/codeation/impress/joint/lazy"
	"github.com/codeation/impress/joint/pipedriver"
	"github.com/codeation/impress/joint/serversocket"
	"github.com/codeation/impress/joint/socketpipes"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

func New() driver.Driver {
	flag.Parse()
	pipeCreator := socketpipes.NewSocketPipes(&socketpipes.Config{
		Listen:      *listen,
		FS:          os.DirFS(*dir),
		NewSocketFn: func() socketpipes.SocketHandler { return serversocket.New() },
	})
	return lazy.New(pipedriver.New(pipeCreator))
}
