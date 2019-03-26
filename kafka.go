package main

import (
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var kafkaConfig = &kafka.ConfigMap{
	"auto.offset.reset": "latest",
}

var kafkaConsumer *kafka.Consumer
var lock = &sync.Mutex{}
var consumerOwners = 0

func getKafkaConsumer() *kafka.Consumer {
	lock.Lock()
	defer lock.Unlock()
	consumerOwners++

	if kafkaConsumer != nil {
		return kafkaConsumer
	}

	c, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		log.Fatalf("Create kafka consumer: %v", err)
	}

	err = c.Subscribe("trackbeat-debug", nil)
	if err != nil {
		log.Fatalf("Subscribe to topic: %v", err)
	}

	kafkaConsumer = c
	return c
}

func closeKafkaConsumer(c *kafka.Consumer) {
	lock.Lock()
	defer lock.Unlock()
	if consumerOwners > 0 {
		consumerOwners--
	}

	if kafkaConsumer != nil && consumerOwners == 0 {
		err := c.Close()
		if err != nil {
			log.Fatalf("Close consumer: %v", err)
		}

		kafkaConsumer = nil
	}
}
