package canvas

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/codeation/impress"
	"github.com/codeation/impress/joint/client"
	"github.com/codeation/impress/joint/domain"
	"github.com/codeation/impress/joint/eventchan"
	"github.com/codeation/impress/joint/iosplit"
	"github.com/codeation/impress/joint/rpc"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

func init() {
	flag.Parse()
	log.Printf("listening on %q...", *listen)

	streamLink := iosplit.NewIOSplit().WithTimeout()
	requestLink := iosplit.NewIOSplit().WithTimeout()
	responseR, responseW := io.Pipe()
	eventR, eventW := io.Pipe()

	eventPipe := rpc.NewPipe(rpc.WithoutMutex(), nil, eventR)
	streamPipe := rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(streamLink), nil)
	syncPipe := rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(requestLink), responseR)

	eventChan := eventchan.New()
	client := client.New(eventChan, eventPipe, streamPipe, syncPipe)
	driver := domain.New(client, eventChan)
	impress.Register(driver)

	go linkRun(streamLink, requestLink, responseW, eventW)

}

type dataStream struct {
	handlers map[string]http.Handler
	readers  map[string]io.Writer
	writers  map[string]io.Reader
}

func newDataIO(streamLink, requestLink io.Reader, responseLink, eventLink io.Writer) *dataStream {
	fileServer := http.FileServer(http.Dir(*dir))
	s := &dataStream{
		handlers: map[string]http.Handler{},
		readers:  map[string]io.Writer{},
		writers:  map[string]io.Reader{},
	}
	for _, path := range []string{"/", "/index.html", "/main.wasm", "/wasm_exec.js"} {
		s.handlers[path] = fileServer
	}
	s.AddWriter("/stream", streamLink)
	s.AddWriter("/request", requestLink)
	s.AddReader("/response", responseLink)
	s.AddReader("/event", eventLink)
	return s
}

func (s *dataStream) AddReader(path string, r io.Writer) {
	s.readers[path] = r
}

func (s *dataStream) AddWriter(path string, w io.Reader) {
	s.writers[path] = w
}

func (s *dataStream) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := s.handlers[req.URL.Path]; ok {
		handler.ServeHTTP(w, req)
		return
	}

	if writer, ok := s.writers[req.URL.Path]; ok {
		w.Header().Add("Content-Type", "application/binary")
		if _, err := io.Copy(w, writer); err != nil {
			log.Printf("io.Copy: %v", err)
		}
		return
	}

	if reader, ok := s.readers[req.URL.Path]; ok {
		if _, err := io.Copy(reader, req.Body); err != nil {
			log.Printf("io.Copy: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func linkRun(streamLink, requestLink io.Reader, responseLink, eventLink io.Writer) error {
	s := &http.Server{
		Addr:           *listen,
		Handler:        newDataIO(streamLink, requestLink, responseLink, eventLink),
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.ListenAndServe()
}
