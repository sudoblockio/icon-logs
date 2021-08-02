package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/geometry-labs/icon-logs/api/healthcheck"
	"github.com/geometry-labs/icon-logs/api/routes"
	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/global"
	"github.com/geometry-labs/icon-logs/kafka"
	"github.com/geometry-labs/icon-logs/logging"
	"github.com/geometry-labs/icon-logs/metrics"
	_ "github.com/geometry-labs/icon-logs/models" // for swagger docs
)

func main() {
	config.ReadEnvironment()

	logging.Init()
	zap.S().Debug("Main: Starting logging with level ", config.Config.LogLevel)

	// Start kafka consumers
	// Go routines start in function
	kafka.StartAPIConsumers()

	// Start Prometheus client
	// Go routine starts in function
	metrics.MetricsAPIStart()

	// Start API server
	// Go routine starts in function
	routes.Start()

	// Start Health server
	// Go routine starts in function
	healthcheck.Start()

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
