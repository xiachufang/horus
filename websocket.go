package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

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
		log.Errorf("Upgrade Error: %v\n", err)
		return
	}
	defer func() {
		conn.Close()
	}()

	clientCloseCh := make(chan struct{})
	conn.SetCloseHandler(func(code int, text string) error {
		message := websocket.FormatCloseMessage(code, "")
		conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		log.Debugf("Client closed connection")
		close(clientCloseCh)
		return nil
	})

	go readLoop(conn)

	log.Debugf("subscribe")
	ch := subscribe()
	defer unsubscribe(ch)

	log.Debugf("begin handler's for loop")
	for {
		select {
		case msg, more := <-ch:
			if !more {
				log.Debugf("close ch in handler")
				return
			}
			// log.Debugf("receive msg in handler")

			value := msg.([]byte)
			if err := conn.WriteMessage(websocket.TextMessage, value); err != nil {
				log.Debugf("Write message Error: %v\n", err)
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
