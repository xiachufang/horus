package main

import (
	"sync"
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

var startKafkaConsumerLock = &sync.Mutex{}
var consumers = 0

func startKafkaConsumer() {
	startKafkaConsumerLock.Lock()
	currentConsumers := consumers
	consumers++
	startKafkaConsumerLock.Unlock()

	if currentConsumers == 0 {
		go func() {
			log.Debugf("Start background kafka consumer\n")
			defer func() {
				stopKafkaConsumer()
				log.Debugf("Finish background kafka consumer\n")
			}()

			kafkaConsumer := createKafkaConsumer()
			defer closeKafkaConsumer(kafkaConsumer)

			running := true
			for running {
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

				startKafkaConsumerLock.Lock()
				running = consumers > 0
				startKafkaConsumerLock.Unlock()
			}
		}()
	} else {
		log.Debugf("Background kafka consumer is running, skip starting\n")
	}
}

func stopKafkaConsumer() {
	log.Debugf("Stop background kafka consumer.")
	startKafkaConsumerLock.Lock()
	defer startKafkaConsumerLock.Unlock()
	if consumers > 0 {
		consumers--
	}
}
