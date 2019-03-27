package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
)

var shutdownCh = make(chan struct{})

func main() {
	var args struct {
		Listen       string
		KafkaBroker  string
		KafkaGroupID string
		KafkaTopic   string
		LogLevel     string `help:"Panic|Fatal|Error|Warn|Info|Debug|Trace"`
	}

	args.Listen = "0.0.0.0:8888"
	args.KafkaBroker = "kafka-01:9092,kafka-02:9092,kafka-03:9092"
	args.KafkaGroupID = "horus"
	args.KafkaTopic = "trackbeat-debug"
	args.LogLevel = "Warn"

	arg.MustParse(&args)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
	var level log.Level
	switch args.LogLevel {
	case "Panic":
		level = log.PanicLevel
	case "Fatal":
		level = log.FatalLevel
	case "Error":
		level = log.ErrorLevel
	case "Warn":
		level = log.WarnLevel
	case "Info":
		level = log.InfoLevel
	case "Debug":
		level = log.DebugLevel
	case "Trace":
		level = log.TraceLevel
	default:
		level = log.WarnLevel
	}
	log.SetLevel(level)

	kafkaConfig.SetKey("bootstrap.servers", args.KafkaBroker)
	kafkaConfig.SetKey("group.id", args.KafkaGroupID)
	kafkaTopic = args.KafkaTopic

	runPubSub()

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
		fmt.Println("\nShutting down the server...")
		close(shutdownCh)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		s.Shutdown(ctx)
		cancel()
	}()

	fmt.Printf("Listening on %s\n", args.Listen)
	s.ListenAndServe()
}
