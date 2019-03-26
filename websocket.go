package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	defer conn.Close()

	go readLoop(conn)

	kafkaConsumer := getKafkaConsumer()
	defer closeKafkaConsumer(kafkaConsumer)

	for {
		msg, err := kafkaConsumer.ReadMessage(timeout)
		if err == nil {
			if err := conn.WriteMessage(websocket.TextMessage, msg.Value); err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			if err.(kafka.Error).Code() != kafka.ErrTimedOut {
				// The client will automatically try to recover from all errors.
				log.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}
}
