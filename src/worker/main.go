package main

import (
	"log"

	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/global"
	"github.com/geometry-labs/icon-logs/kafka"
	"github.com/geometry-labs/icon-logs/logging"
	"github.com/geometry-labs/icon-logs/metrics"
	"github.com/geometry-labs/icon-logs/worker/routines"
	"github.com/geometry-labs/icon-logs/worker/transformers"
)

func main() {
	config.ReadEnvironment()

	logging.Init()
	log.Printf("Main: Starting logging with level %s", config.Config.LogLevel)

	// Start Prometheus client
	metrics.Start()

	// Feature flags
	if config.Config.OnlyRunAllRoutines == true {
		// Start routines
		routines.StartLogCountRoutine()
		routines.StartLogCountByAddressRoutine()
		// routines.StartLogMissingRoutine()

		global.WaitShutdownSig()
	}

	// Start kafka consumer
	kafka.StartWorkerConsumers()

	// Start transformers
	transformers.StartLogsTransformer()

	global.WaitShutdownSig()
}
