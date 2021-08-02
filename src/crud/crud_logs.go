package crud

import (
	"sync"
	"context"

	"go.uber.org/zap"
	"github.com/cenkalti/backoff/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/models"
)

type LogModel struct {
	mongoConn        *MongoConn
	WriteChan        chan *models.Log
}

var logModelInstance *LogModel
var logModelOnce sync.Once

func GetLogModel() *LogModel {
	logModelOnce.Do(func() {
		logModelInstance = &LogModel{
			mongoConn:        GetMongoConn(),
			WriteChan:        make(chan *models.Log, 1),
		}

		// logModelInstance.CreateNumberIndex("blocknumber", false, false)
		// logModelInstance.CreateStringIndex("toaddress")
	})
	return logModelInstance
}

func (b *LogModel) getCollectionHandle() *mongo.Collection {
  dbName := config.Config.DbName
  dbCollection := config.Config.DbCollection

  return GetMongoConn().DatabaseHandle(dbName).Collection(dbCollection)
}

func (b *LogModel) CreateNumberIndex(field string, isAscending bool, isUnique bool) {
	ascending := 1
	if !isAscending {
		ascending = -1
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.M{field: ascending},
		Options: options.Index().SetUnique(isUnique),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := b.getCollectionHandle().Indexes().CreateOne(ctx, indexModel)
  if err != nil {
    zap.S().Fatal("CREATENUMBERINDEX PANIC: ", err.Error())
  }
}

func (b *LogModel) CreateStringIndex(field string) {

	indexModel := mongo.IndexModel{
		Keys:    bson.M{field: "hashed"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := b.getCollectionHandle().Indexes().CreateOne(ctx, indexModel)
  if err != nil {
    zap.S().Fatal("CREATESTRINGINDEX PANIC: ", err.Error())
  }
}

func (b *LogModel) Insert(ctx context.Context, log *models.Log) error {

  err := backoff.Retry(func() error {
	  _, err := b.getCollectionHandle().InsertOne(ctx, log)

		if err != nil {
			zap.S().Info("MongoDb RetryCreate Error : ", err.Error())
		}

    return err
	}, backoff.NewExponentialBackOff())

	return err
}

func (b *LogModel) Select(
	ctx context.Context,
	limit int64,
	skip int64,
	hash string,
	from string,
	to string,
) ([]models.Log, error) {
  err := b.mongoConn.retryPing(ctx)
  if err != nil {
    return nil, err
  }

	if limit <= 0 {
		limit = 1
	} else if limit > 100 {
		limit = 100
	}
	if skip < 0 {
		skip = 0
	}

	// Building KeyValue pairs
	queryParams := make(map[string]interface{})
	// hash
	if hash != "" {
		queryParams["hash"] = hash
	}
	// from
	if from != "" {
		queryParams["fromaddress"] = from
	}
	// to
	if to != "" {
		queryParams["toaddress"] = to
	}

	// Building FindOptions
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
  opts.SetSort(bson.D{{"blocknumber", -1}})

	queryParamsD, err := convertMapToBsonD(queryParams)
	if err != nil {
    return nil, err
	}

	cursor, err := b.getCollectionHandle().Find(ctx, queryParamsD, &opts)
	if err != nil {
    return nil, err
	}

	var results []bson.M
  err = cursor.All(ctx, &results)
	if err != nil {
    return nil, err
	}

  // convert bson to model
  logs := make([]models.Log, 0)
  for _, r := range results {
    logs = append(logs, convertBsonMToLog(r))
  }

	return logs, nil
}

func StartLogLoader() {

  go func() {
    var log *models.Log
    mongoLoaderChan := GetLogModel().WriteChan

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    for {
      log = <-mongoLoaderChan
      GetLogModel().Insert(ctx, log)

      zap.S().Info("Loader: Loaded in collection Logs - BlockNumber=", log.BlockNumber)
    }
  }()
}
