package routines

import (
	"time"

	"go.uber.org/zap"

	"github.com/geometry-labs/icon-logs/crud"
	"github.com/geometry-labs/icon-logs/models"
	"github.com/geometry-labs/icon-logs/worker/utils"
)

func StartLogMissingRoutine() {

	// routine every day
	go logMissingRoutine(3600 * time.Second)
}

func logMissingRoutine(duration time.Duration) {

	// Loop every duration
	for {

		currentBlockNumber := 1

		for {

			transactionHashes, err := utils.IconNodeServiceGetBlockTransactionHashes(currentBlockNumber)
			if err != nil {
				zap.S().Warn(
					"Routine=LogMissing",
					" CurrentBlockNumber=", currentBlockNumber,
					" Error=", err.Error(),
					" Sleeping 1 second...",
				)

				time.Sleep(1 * time.Second)
				continue
			}

			for _, txHash := range *transactionHashes {

				// loop until success
				for {
					logCount, err := utils.IconNodeServiceGetTransactionLogCount(txHash)

					if err != nil {
						zap.S().Warn(
							"Routine=LogMissing",
							" CurrentBlockNumber=", currentBlockNumber,
							" TransactionHash=", txHash,
							" Error=", err.Error(),
							" Sleeping 1 second...",
						)

						time.Sleep(1 * time.Second)
						continue
					}

					logs, err := crud.GetLogModel().SelectMany(3000, 0, 0, txHash, "", "")

					if len(*logs) != logCount {

						logMissing := &models.LogMissing{
							TransactionHash: txHash,
						}

						crud.GetLogMissingModel().LoaderChannel <- logMissing
					}

					break
				}
			}

			if currentBlockNumber > 44000000 {
				break
			}

			if currentBlockNumber%100000 == 0 {
				zap.S().Info("Routine=LogMissing, CurrentBlockNumber= ", currentBlockNumber, " - Checked 100,000 blocks...")
			}

			currentBlockNumber++
		}

		zap.S().Info("Completed routine, sleeping...")
		time.Sleep(duration)
	}
}
