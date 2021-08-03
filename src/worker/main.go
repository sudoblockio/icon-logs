package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/crud"
	"github.com/geometry-labs/icon-logs/global"
	"github.com/geometry-labs/icon-logs/kafka"
	"github.com/geometry-labs/icon-logs/logging"
	"github.com/geometry-labs/icon-logs/metrics"
	"github.com/geometry-labs/icon-logs/worker/transformers"
)

func main() {
	config.ReadEnvironment()

	logging.Init()
	log.Printf("Main: Starting logging with level %s", config.Config.LogLevel)

	// Start Prometheus client
	metrics.MetricsWorkerStart()

	// Start Mongodb loader
	crud.StartLogLoader()

	// Start kafka consumer
	kafka.StartWorkerConsumers()

	// Start kafka Producer
	kafka.StartProducers()

	// Start transformers
	transformers.StartLogsTransformer()

	// Listen for close sig
	// Register for interupt (Ctrl+C) and SIGTERM (docker)

	//create a notification channel to shutdown
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		zap.S().Info("Shutting down...")
		global.ShutdownChan <- 1
	}()

	<-global.ShutdownChan
}
