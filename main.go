package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			break
		}
	}
}

// AppHandler is for handers with context
type AppHandler struct {
	events chan []byte
}

// NewAppHandler create new instance
func NewAppHandler() *AppHandler {
	return &AppHandler{
		events: make(chan []byte),
	}
}

func (app *AppHandler) sendEvent(e []byte) {
	app.events <- e
}

func (app *AppHandler) close() {
	close(app.events)
}

func (app *AppHandler) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	go readLoop(conn)

	for {
		event, more := <-app.events
		if !more {
			return
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, event); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	app := NewAppHandler()
	http.HandleFunc("/topic", app.handler)
	go func() {
		log.Println("Listening on :8881")
		err := http.ListenAndServe("0.0.0.0:8881", nil)
		if err != nil {
			log.Fatalf("listen: %v", err)
		}
	}()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka-01:9092,kafka-02:9092,kafka-03:9092",
		"group.id":          "horus",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		log.Fatalf("Create kafka consumer: %v", err)
	}

	err = c.Subscribe("trackbeat-debug", nil)
	if err != nil {
		log.Fatalf("Subscribe to topic: %v", err)
	}

	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Signal %s received, exiting..", sig)
			run = false
		default:
			msg, err := c.ReadMessage(3)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				app.sendEvent(msg.Value)
			} else if err.(kafka.Error).Code() == kafka.ErrTimedOut {

			} else {
				// The client will automatically try to recover from all errors.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}

	app.close()
	c.Close()
}
