package transformers

import (
	"encoding/hex"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/crud"
	"github.com/geometry-labs/icon-logs/kafka"
	"github.com/geometry-labs/icon-logs/models"
)

func StartLogsTransformer() {
	go logsTransformer()
}

func logsTransformer() {
	consumerTopicNameLogs := config.Config.ConsumerTopicLogs

	// Input Channels
	consumerTopicChanLogs := kafka.KafkaTopicConsumers[consumerTopicNameLogs].TopicChan

	// Output channels
	logLoaderChan := crud.GetLogModel().WriteChan
	logCountLoaderChan := crud.GetLogCountModel().WriteChan

	zap.S().Debug("Logs Worker: started working")
	for {
		// Read from kafka
		consumerTopicMsg := <-consumerTopicChanLogs

		// Log message from ETL
		logRaw, err := convertBytesToLogRawProtoBuf(consumerTopicMsg.Value)
		if err != nil {
			zap.S().Fatal("Logs Worker: Unable to proceed cannot convert kafka msg value to LogRaw, err: ", err.Error())
		}

		// Transform logic
		var transformedLog *models.Log
		transformedLog = transformLogRaw(logRaw)

		// Load log to Postgres
		logLoaderChan <- transformedLog

		// Load log counter to Postgres
		logCount := &models.LogCount{
			Count: 1, // Adds with current
			Id:    1, // Only one row
		}
		logCountLoaderChan <- logCount

		zap.S().Debug("Logs worker: last seen log #", string(consumerTopicMsg.Key))
	}
}

func convertBytesToLogRawProtoBuf(value []byte) (*models.LogRaw, error) {
	log := models.LogRaw{}
	err := proto.Unmarshal(value[6:], &log)
	if err != nil {
		zap.S().Error("Error: ", err.Error())
		zap.S().Error("Value=", hex.Dump(value[6:]))
	}
	return &log, err
}

// Business logic goes here
func transformLogRaw(logRaw *models.LogRaw) *models.Log {
	return &models.Log{
		Type:             logRaw.Type,
		LogIndex:         logRaw.LogIndex,
		TransactionHash:  logRaw.TransactionHash,
		TransactionIndex: logRaw.TransactionIndex,
		Address:          logRaw.Address,
		Data:             logRaw.Data,
		Indexed:          logRaw.Indexed,
		BlockNumber:      logRaw.BlockNumber,
		BlockTimestamp:   logRaw.BlockTimestamp,
		BlockHash:        logRaw.BlockHash,
		ItemId:           logRaw.ItemId,
		ItemTimestamp:    logRaw.ItemTimestamp,
	}
}
