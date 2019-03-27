package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

const timeout = 3

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			break
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer log.Println("Finish handler")
	defer conn.Close()

	clientCloseCh := make(chan struct{})
	conn.SetCloseHandler(func(code int, text string) error {
		log.Println("Client closed connection")
		close(clientCloseCh)
		return nil
	})

	go readLoop(conn)

	ch := subscribe()
	defer unsubscribe(ch)

	startKafkaConsumer()
	defer stopKafkaConsumer()

	for {
		select {
		case msg, more := <-ch:
			if !more {
				return
			}

			value := msg.([]byte)
			if err := conn.WriteMessage(websocket.TextMessage, value); err != nil {
				log.Println(err)
				return
			}
		case <-shutdownCh:
			return
		case <-clientCloseCh:
			return
		}
	}
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

// Server is websocket server
type Server struct {
	mux *http.ServeMux
}

func newServer() *Server {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		mux: http.NewServeMux(),
	}

	s.mux.HandleFunc("/topic", handler)

	fs := http.FileServer(&neuteredFileSystem{http.Dir(dir)})
	s.mux.Handle("/", fs)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
