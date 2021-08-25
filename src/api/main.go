package main

import (
	"log"

	"github.com/geometry-labs/icon-logs/api/healthcheck"
	"github.com/geometry-labs/icon-logs/api/routes"
	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/global"
	"github.com/geometry-labs/icon-logs/logging"
	"github.com/geometry-labs/icon-logs/metrics"
	_ "github.com/geometry-labs/icon-logs/models" // for swagger docs
)

func main() {
	config.ReadEnvironment()

	logging.Init()
	log.Printf("Main: Starting logging with level %s", config.Config.LogLevel)

	// Start Prometheus client
	// Go routine starts in function
	metrics.MetricsAPIStart()

	// Start API server
	// Go routine starts in function
	routes.Start()

	// Start Health server
	// Go routine starts in function
	healthcheck.Start()

	global.WaitShutdownSig()
}
