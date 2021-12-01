package routines

import (
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/geometry-labs/icon-logs/crud"
	"github.com/geometry-labs/icon-logs/models"
	"github.com/geometry-labs/icon-logs/redis"
)

func StartLogCountByAddressRoutine() {

	// routine every day
	go logCountByAddressRoutine(3600 * time.Second)
}

func logCountByAddressRoutine(duration time.Duration) {

	// Loop every duration
	for {

		// Loop through all addresses
		skip := 0
		limit := 100
		for {
			addresses, err := crud.GetLogCountByAddressModel().SelectMany(limit, skip)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Sleep
				zap.S().Info("Routine=LogCountByAddress", " - No records found, sleeping...")
				break
			} else if err != nil {
				zap.S().Fatal(err.Error())
			}
			if len(*addresses) == 0 {
				// Sleep
				break
			}

			zap.S().Info("Routine=LogCountByAddress", " - Processing ", len(*addresses), " addresses...")
			for _, a := range *addresses {

				///////////
				// Count //
				///////////
				count, err := crud.GetLogModel().CountByAddress(a.Address)
				if err != nil {
					// Postgres error
					zap.S().Warn(err)
					continue
				}

				//////////////////
				// Update Redis //
				//////////////////
				countKey := "icon_logs_log_count_by_address_" + a.Address
				err = redis.GetRedisClient().SetCount(countKey, count)
				if err != nil {
					// Redis error
					zap.S().Warn(err)
					continue
				}

				/////////////////////
				// Update Postgres //
				/////////////////////
				logCountByAddress := &models.LogCountByAddress{
					Address: a.Address,
					Count:   uint64(count),
				}
				err = crud.GetLogCountByAddressModel().UpsertOne(logCountByAddress)
			}

			skip += limit
		}

		zap.S().Info("Completed routine, sleeping...")
		time.Sleep(duration)
	}
}
