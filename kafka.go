package main

import (
	"log"
	"sync"

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

var startKafkaConsumerLock = &sync.Mutex{}
var consumers = 0

func startKafkaConsumer() {
	startKafkaConsumerLock.Lock()
	defer startKafkaConsumerLock.Unlock()

	if consumers == 0 {
		go func() {
			log.Println("Start background kafka consumer")
			defer func() {
				stopKafkaConsumer()
				log.Println("Finish background kafka consumer")
			}()

			kafkaConsumer := createKafkaConsumer()
			defer closeKafkaConsumer(kafkaConsumer)

			running := true
			for running {
				msg, err := kafkaConsumer.ReadMessage(timeout)
				if err == nil {
					log.Printf("Publish Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
					if !publish(msg.Value) {
						log.Printf("No subscribers now\n")
						return
					}
					log.Printf("Finish Publish\n")
				} else {
					if err.(kafka.Error).Code() != kafka.ErrTimedOut {
						// The client will automatically try to recover from all errors.
						log.Printf("Consumer error: %v (%v)\n", err, msg)
					}
				}

				startKafkaConsumerLock.Lock()
				running = consumers > 0
				startKafkaConsumerLock.Unlock()
			}
		}()
	} else {
		log.Printf("Background kafka consumer is running, skip starting\n")
	}

	consumers++
}

func stopKafkaConsumer() {
	log.Printf("Stop background kafka consumer.")
	startKafkaConsumerLock.Lock()
	defer startKafkaConsumerLock.Unlock()
	if consumers > 0 {
		consumers--
	}
}
