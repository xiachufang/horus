package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var kafkaConfig = &kafka.ConfigMap{
	"auto.offset.reset": "latest",
}
var kafkaTopic = "trackbeat-debug"

func createKafkaConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		log.Fatalf("Create kafka consumer: %v", err)
	}

	err = c.Subscribe(kafkaTopic, nil)
	if err != nil {
		log.Fatalf("Subscribe to topic: %v", err)
	}

	return c
}

func closeKafkaConsumer(c *kafka.Consumer) {
	if c != nil {
		err := c.Close()
		if err != nil {
			log.Fatalf("Close consumer: %v", err)
		}
	}
}

func startKafkaConsumer() {
	go func() {
		log.Debugf("Start background kafka consumer\n")

		kafkaConsumer := createKafkaConsumer()
		defer closeKafkaConsumer(kafkaConsumer)

		for {
			select {
			case <-shutdownCh:
				return
			default:
				log.Debug("Read message")
				msg, err := kafkaConsumer.ReadMessage(timeout * time.Second)
				if err == nil {
					log.Debugf("Publish Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
					publish(msg.Value)
					log.Debugf("Finish Publish\n")
				} else {
					if err.(kafka.Error).Code() != kafka.ErrTimedOut {
						// The client will automatically try to recover from all errors.
						log.Errorf("Consumer error: %v (%v)\n", err, msg)
					} else {
						log.Debugf("Read message timeout")
					}
				}
			}
		}
	}()
}
