package main

import (
	"log"
	"net/http"

	"github.com/alexflint/go-arg"
)

func main() {

	var args struct {
		Listen       string
		KafkaBroker  string
		KafkaGroupID string
	}

	args.Listen = "0.0.0.0:8888"
	args.KafkaBroker = "kafka-01:9092,kafka-02:9092,kafka-03:9092"
	args.KafkaGroupID = "horus"

	arg.MustParse(&args)
	kafkaConfig.SetKey("bootstrap.servers", args.KafkaBroker)
	kafkaConfig.SetKey("group.id", args.KafkaGroupID)
	http.HandleFunc("/topic", handler)
	log.Printf("Listening on %s", args.Listen)
	http.ListenAndServe(args.Listen, nil)
}
