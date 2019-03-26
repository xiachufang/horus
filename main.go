package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
)

var shutdownCh = make(chan struct{})

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

	s := &http.Server{
		Addr:           args.Listen,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        newServer(),
	}

	go func() {
		var sigchan = make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Println("\nShutting down the server...")
		close(shutdownCh)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		s.Shutdown(ctx)
		cancel()
	}()

	kafkaConfig.SetKey("bootstrap.servers", args.KafkaBroker)
	kafkaConfig.SetKey("group.id", args.KafkaGroupID)

	log.Printf("Listening on %s", args.Listen)
	s.ListenAndServe()
}
