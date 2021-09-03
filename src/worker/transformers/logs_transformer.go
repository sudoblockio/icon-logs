package transformers

import (
	"encoding/hex"
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/crud"
	"github.com/geometry-labs/icon-logs/kafka"
	"github.com/geometry-labs/icon-logs/models"
	"github.com/geometry-labs/icon-logs/redis"
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
	redisClient := redis.GetRedisClient()

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
		log := transformLogRawToLog(logRaw)

		// Push to redis
		logWebsocket := transformLogToLogWS(log)
		logWebsocketJSON, _ := json.Marshal(logWebsocket)
		redisClient.Publish(logWebsocketJSON)

		// Load log to Postgres
		logLoaderChan <- log

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
func transformLogRawToLog(logRaw *models.LogRaw) *models.Log {
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

// Business logic goes here
func transformLogToLogWS(log *models.Log) *models.LogWebsocket {
	return &models.LogWebsocket{
		Type:             log.Type,
		LogIndex:         log.LogIndex,
		TransactionHash:  log.TransactionHash,
		TransactionIndex: log.TransactionIndex,
		Address:          log.Address,
		Data:             log.Data,
		Indexed:          log.Indexed,
		BlockNumber:      log.BlockNumber,
		BlockTimestamp:   log.BlockTimestamp,
		BlockHash:        log.BlockHash,
	}
}
