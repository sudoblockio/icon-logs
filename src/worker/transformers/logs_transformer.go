package transformers

import (
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
  consumer_topic_name_logs := "logs"
	producer_topic_name := "logs-ws"

	// Check topic names
	if utils.StringInSlice(consumer_topic_name_logs, config.Config.ConsumerTopics) == false {
		zap.S().Panic("Logs Worker: no ", consumer_topic_name_logs, " topic found in CONSUMER_TOPICS=", config.Config.ConsumerTopics)
	}
	if utils.StringInSlice(producer_topic_name, config.Config.ProducerTopics) == false {
		zap.S().Panic("Logs Worker: no ", producer_topic_name, " topic found in PRODUCER_TOPICS=", config.Config.ConsumerTopics)
	}

	consumer_topic_chan_logs := make(chan *sarama.ConsumerMessage)
	producer_topic_chan := kafka.KafkaTopicProducers[producer_topic_name].TopicChan
	mongoLoaderChan := crud.GetLogModel().WriteChan

	// Register consumer channel logs
	broadcaster_output_chan_id_log := kafka.Broadcasters[consumer_topic_name_logs].AddBroadcastChannel(consumer_topic_chan_logs)
	defer func() {
		kafka.Broadcasters[consumer_topic_name_logs].RemoveBroadcastChannel(broadcaster_output_chan_id_log)
	}()

	zap.S().Debug("Logs Worker: started working")
	for {
		// Read from kafka
    var consumer_topic_msg *sarama.ConsumerMessage
    var transformedLog *models.Log

    // Log message from ETL
    logRaw, err := convertBytesToLogRaw(consumer_topic_msg.Value)
    if err != nil {
      zap.S().Fatal("Logs Worker: Unable to proceed cannot convert kafka msg value to LogRaw, err: ", err.Error())
    }

    // Transform logic
    transformedLog = transformLogRaw(logRaw)

		// Produce to Kafka
		producer_topic_msg := &sarama.ProducerMessage{
			Topic: producer_topic_name,
			Key:   sarama.ByteEncoder(consumer_topic_msg.Key),
			Value: sarama.ByteEncoder(consumer_topic_msg.Value),
		}
		producer_topic_chan <- producer_topic_msg

		// Load to Postgres
		mongoLoaderChan <- transformedLog

		zap.S().Debug("Logs worker: last seen log #", string(consumer_topic_msg.Key))
	}
}

func convertBytesToLogRaw(value []byte) (*models.LogRaw, error) {
	tx := models.LogRaw{}

	err := protojson.Unmarshal(value, &tx)
	if err != nil {
    zap.S().Panic("Error: ", err.Error(), " Value: ", string(value))
	}

	return &tx, nil
}

func convertBytesToLogRaw(value []byte) (*models.LogRaw, error) {
	log := models.LogRaw{}

	err := protojson.Unmarshal(value, &log)
	if err != nil {
    zap.S().Panic("Error: ", err.Error(), " Value: ", string(value))
	}

	return &log, nil
}

// Business logic goes here
func transformLogRaw(logRaw *models.LogRaw) *models.Log {
  return &models.Log {
    Type: logRaw.Type,
    LogIndex: logRaw.LogIndex,
    TransactionHash: logRaw.TransactionHash,
    TransactionIndex: logRaw.TransactionIndex,
    Address: logRaw.Address,
    Data: logRaw.Data,
    Indexed: logRaw.Indexed,
    BlockNumber: logRaw.BlockNumber,
    BlockTimestamp: logRaw.BlockTimestamp,
    BlockHash: logRaw.BlockHash,
    ItemId: logRaw.ItemId,
    ItemTimestamp: logRaw.ItemTimestamp,
  }
}
