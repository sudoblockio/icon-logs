package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type configType struct {
	Name        string `envconfig:"NAME" required:"false" default:"logs-service"`
	NetworkName string `envconfig:"NETWORK_NAME" required:"false" default:"mainnnet"`

	// Ports
	Port        string `envconfig:"PORT" required:"false" default:"8000"`
	HealthPort  string `envconfig:"HEALTH_PORT" required:"false" default:"8180"`
	MetricsPort string `envconfig:"METRICS_PORT" required:"false" default:"9400"`

	// Prefix
	RestPrefix      string `envconfig:"REST_PREFIX" required:"false" default:"/api/v1"`
	WebsocketPrefix string `envconfig:"WEBSOCKET_PREFIX" required:"false" default:"/ws/v1"`
	HealthPrefix    string `envconfig:"HEALTH_PREFIX" required:"false" default:"/health"`
	MetricsPrefix   string `envconfig:"METRICS_PREFIX" required:"false" default:"/metrics"`

	// Endpoints
	MaxPageSize int `envconfig:"MAX_PAGE_SIZE" required:"false" default:"100"`
	MaxPageSkip int `envconfig:"MAX_PAGE_SKIP" required:"false" default:"1000000"`

	// Icon node service
	IconNodeServiceURL string `envconfig:"ICON_NODE_SERVICE_URL" required:"false" default:"https://ctz.solidwallet.io/api/v3"`

	// CORS
	CORSAllowOrigins  string `envconfig:"CORS_ALLOW_ORIGINS" required:"false" default:"*"`
	CORSAllowHeaders  string `envconfig:"CORS_ALLOW_HEADERS" required:"false" default:"*"`
	CORSAllowMethods  string `envconfig:"CORS_ALLOW_METHODS" required:"false" default:"GET,POST,HEAD,PUT,DELETE,PATCH"`
	CORSExposeHeaders string `envconfig:"CORS_EXPOSE_HEADERS" required:"false" default:"*"`

	// Compress
	RestCompressLevel int `envconfig:"REST_COMPRESS_LEVEL" required:"false" default:"2"`

	// Monitoring
	HealthPollingInterval int `envconfig:"HEALTH_POLLING_INTERVAL" required:"false" default:"10"`

	// Logging
	LogLevel         string `envconfig:"LOG_LEVEL" required:"false" default:"INFO"`
	LogToFile        bool   `envconfig:"LOG_TO_FILE" required:"false" default:"false"`
	LogFileName      string `envconfig:"LOG_FILE_NAME" required:"false" default:"logs-service.log"`
	LogFormat        string `envconfig:"LOG_FORMAT" required:"false" default:"json"`
	LogIsDevelopment bool   `envconfig:"LOG_IS_DEVELOPMENT" required:"false" default:"true"`

	// Kafka
	KafkaBrokerURL    string `envconfig:"KAFKA_BROKER_URL" required:"false" default:"localhost:29092"`
	SchemaRegistryURL string `envconfig:"SCHEMA_REGISTRY_URL" required:"false" default:"localhost:8081"`
	KafkaGroupID      string `envconfig:"KAFKA_GROUP_ID" required:"false" default:"logs-service"`

	// Topics
	ConsumerGroup                string            `envconfig:"CONSUMER_GROUP" required:"false" default:"logs-consumer-group"`
	ConsumerIsTail               bool              `envconfig:"CONSUMER_IS_TAIL" required:"false" default:"false"`
	ConsumerJobID                string            `envconfig:"CONSUMER_JOB_ID" required:"false" default:""`
	ConsumerGroupBalanceStrategy string            `envconfig:"CONSUMER_GROUP_BALANCE_STRATEGY" required:"false" default:"BalanceStrategySticky"`
	ConsumerTopicBlocks          string            `envconfig:"CONSUMER_TOPIC_BLOCKS" required:"false" default:"blocks"`
	ConsumerTopicTransactions    string            `envconfig:"CONSUMER_TOPIC_TRANSACTIONS" required:"false" default:"transactions"`
	ConsumerTopicLogs            string            `envconfig:"CONSUMER_TOPIC_LOGS" required:"false" default:"logs"`
	ConsumerIsPartitionConsumer  bool              `envconfig:"CONSUMER_IS_PARTITION_CONSUMER" required:"false" default:"false"`
	ConsumerPartition            int               `envconfig:"CONSUMER_PARTITION" required:"false" default:"0"`
	ConsumerPartitionTopic       string            `envconfig:"CONSUMER_PARTITION_TOPIC" required:"false" default:"logs"`
	ConsumerPartitionStartOffset int               `envconfig:"CONSUMER_PARTITION_START_OFFSET" required:"false" default:"1"`
	ProducerTopics               []string          `envconfig:"PRODUCER_TOPICS" required:"false" default:"logs-ws"`
	SchemaNameTopics             map[string]string `envconfig:"SCHEMA_NAME_TOPICS" required:"false" default:"logs-ws:logs"`
	SchemaFolderPath             string            `envconfig:"SCHEMA_FOLDER_PATH" required:"false" default:"schemas/"`

	// DB
	DbDriver             string `envconfig:"DB_DRIVER" required:"false" default:"postgres"`
	DbHost               string `envconfig:"DB_HOST" required:"false" default:"localhost"`
	DbPort               string `envconfig:"DB_PORT" required:"false" default:"5432"`
	DbUser               string `envconfig:"DB_USER" required:"false" default:"postgres"`
	DbPassword           string `envconfig:"DB_PASSWORD" required:"false" default:"changeme"`
	DbName               string `envconfig:"DB_DBNAME" required:"false" default:"postgres"`
	DbSslmode            string `envconfig:"DB_SSL_MODE" required:"false" default:"disable"`
	DbTimezone           string `envconfig:"DB_TIMEZONE" required:"false" default:"UTC"`
	DbMaxIdleConnections int    `envconfig:"DB_MAX_IDLE_CONNECTIONS" required:"false" default:"2"`
	DbMaxOpenConnections int    `envconfig:"DB_MAX_OPEN_CONNECTIONS" required:"false" default:"10"`

	// Redis
	RedisHost                     string `envconfig:"REDIS_HOST" required:"false" default:"localhost"`
	RedisPort                     string `envconfig:"REDIS_PORT" required:"false" default:"6379"`
	RedisPassword                 string `envconfig:"REDIS_PASSWORD" required:"false" default:""`
	RedisChannel                  string `envconfig:"REDIS_CHANNEL" required:"false" default:"logs"`
	RedisSentinelClientMode       bool   `envconfig:"REDIS_SENTINEL_CLIENT_MODE" required:"false" default:"false"`
	RedisSentinelClientMasterName string `envconfig:"REDIS_SENTINEL_CLIENT_MASTER_NAME" required:"false" default:"master"`

	// GORM
	GormLoggingThresholdMilli int `envconfig:"GORM_LOGGING_THRESHOLD_MILLI" required:"false" default:"250"`

	// Feature flags
	OnlyRunAllRoutines bool `envconfig:"ONLY_RUN_ALL_ROUTINES" required:"false" default:"false"`
}

var Config configType

func ReadEnvironment() {
	err := envconfig.Process("", &Config)
	if err != nil {
		log.Fatalf("ERROR: envconfig - %s\n", err.Error())
	}
}
