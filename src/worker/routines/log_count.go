package routines

import (
	"time"

	"go.uber.org/zap"

	"github.com/geometry-labs/icon-logs/crud"
	"github.com/geometry-labs/icon-logs/models"
	"github.com/geometry-labs/icon-logs/redis"
)

func StartLogCountRoutine() {

	// routine every day
	go logCountRoutine(3600 * time.Second)
}

func logCountRoutine(duration time.Duration) {

	// Loop every duration
	for {

		/////////////
		// Regular //
		/////////////

		// Count
		count, err := crud.GetLogCountIndexModel().Count()
		if err != nil {
			// Postgres error
			zap.S().Warn(err)
			continue
		}

		// Update Redis
		countKey := "icon_logs_log_count_log"
		err = redis.GetRedisClient().SetCount(countKey, count)
		if err != nil {
			// Redis error
			zap.S().Warn(err)
			continue
		}

		// Update Postgres
		logCount := &models.LogCount{
			Type:  "log",
			Count: uint64(count),
		}
		err = crud.GetLogCountModel().UpsertOne(logCount)

		zap.S().Info("Completed routine, sleeping...")
		time.Sleep(duration)
	}
}
