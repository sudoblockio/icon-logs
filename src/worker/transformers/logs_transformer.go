package transformers

import (
	"encoding/hex"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/Shopify/sarama.v1"

	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/crud"
	"github.com/geometry-labs/icon-logs/kafka"
	"github.com/geometry-labs/icon-logs/models"
	"github.com/geometry-labs/icon-logs/worker/utils"
)

func StartLogsTransformer() {
	go logsTransformer()
}

func logsTransformer() {
	consumer_topic_name := "logs"
	producer_topic_name := "logs-ws"

	// Check topic names
	if utils.StringInSlice(consumer_topic_name, config.Config.ConsumerTopics) == false {
		zap.S().Panic("Logs Worker: no ", consumer_topic_name, " topic found in CONSUMER_TOPICS=", config.Config.ConsumerTopics)
	}
	if utils.StringInSlice(producer_topic_name, config.Config.ProducerTopics) == false {
		zap.S().Panic("Logs Worker: no ", producer_topic_name, " topic found in PRODUCER_TOPICS=", config.Config.ConsumerTopics)
	}

	consumer_topic_chan := make(chan *sarama.ConsumerMessage)
	producer_topic_chan := kafka.KafkaTopicProducers[producer_topic_name].TopicChan
	logLoaderChan := crud.GetLogModel().WriteChan
	logCountLoaderChan := crud.GetLogCountModel().WriteChan

	// Register consumer channel logs
	broadcaster_output_chan_id_log := kafka.Broadcasters[consumer_topic_name].AddBroadcastChannel(consumer_topic_chan)
	defer func() {
		kafka.Broadcasters[consumer_topic_name].RemoveBroadcastChannel(broadcaster_output_chan_id_log)
	}()

	zap.S().Debug("Logs Worker: started working")
	for {
		// Read from kafka
		consumer_topic_msg := <-consumer_topic_chan

		// Log message from ETL
		logRaw, err := convertBytesToLogRawProtoBuf(consumer_topic_msg.Value)
		if err != nil {
			zap.S().Fatal("Logs Worker: Unable to proceed cannot convert kafka msg value to LogRaw, err: ", err.Error())
		}

		// Transform logic
		var transformedLog *models.Log
		transformedLog = transformLogRaw(logRaw)

		// Produce to Kafka
		producer_topic_msg := &sarama.ProducerMessage{
			Topic: producer_topic_name,
			Key:   sarama.ByteEncoder(consumer_topic_msg.Key),
			Value: sarama.ByteEncoder(consumer_topic_msg.Value),
		}
		producer_topic_chan <- producer_topic_msg

		// Load log to Postgres
		logLoaderChan <- transformedLog

		// Load log counter to Postgres
		logCount := &models.LogCount{
			Count: 1, // Adds with current
			Id:    1, // Only one row
		}
		logCountLoaderChan <- logCount

		zap.S().Debug("Logs worker: last seen log #", string(consumer_topic_msg.Key))
	}
}

func convertBytesToLogRawJSON(value []byte) (*models.LogRaw, error) {
	log := models.LogRaw{}

	err := protojson.Unmarshal(value, &log)
	if err != nil {
		zap.S().Panic("Error: ", err.Error(), " Value: ", string(value))
	}

	return &log, nil
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
