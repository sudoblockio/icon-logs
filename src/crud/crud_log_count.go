package crud

import (
	"errors"
	"reflect"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/geometry-labs/icon-logs/models"
	"github.com/geometry-labs/icon-logs/redis"
)

// LogCountModel - type for address table model
type LogCountModel struct {
	db            *gorm.DB
	model         *models.LogCount
	modelORM      *models.LogCountORM
	LoaderChannel chan *models.LogCount
}

var logCountModel *LogCountModel
var logCountModelOnce sync.Once

// GetAddressModel - create and/or return the addresss table model
func GetLogCountModel() *LogCountModel {
	logCountModelOnce.Do(func() {
		dbConn := getPostgresConn()
		if dbConn == nil {
			zap.S().Fatal("Cannot connect to postgres database")
		}

		logCountModel = &LogCountModel{
			db:            dbConn,
			model:         &models.LogCount{},
			LoaderChannel: make(chan *models.LogCount, 1),
		}

		err := logCountModel.Migrate()
		if err != nil {
			zap.S().Fatal("LogCountModel: Unable migrate postgres table: ", err.Error())
		}

		StartLogCountLoader()
	})

	return logCountModel
}

// Migrate - migrate logCounts table
func (m *LogCountModel) Migrate() error {
	// Only using LogCountRawORM (ORM version of the proto generated struct) to create the TABLE
	err := m.db.AutoMigrate(m.modelORM) // Migration and Index creation
	return err
}

// Select - select from logCounts table
func (m *LogCountModel) SelectOne(_type string) (*models.LogCount, error) {
	db := m.db

	// Set table
	db = db.Model(&models.LogCount{})

	// Address
	db = db.Where("type = ?", _type)

	logCount := &models.LogCount{}
	db = db.First(logCount)

	return logCount, db.Error
}

// Select - select from logCounts table
func (m *LogCountModel) SelectCount(_type string) (uint64, error) {
	db := m.db

	// Set table
	db = db.Model(&models.LogCount{})

	// Address
	db = db.Where("type = ?", _type)

	logCount := &models.LogCount{}
	db = db.First(logCount)

	count := uint64(0)
	if logCount != nil {
		count = logCount.Count
	}

	return count, db.Error
}

func (m *LogCountModel) UpsertOne(
	logCount *models.LogCount,
) error {
	db := m.db

	// map[string]interface{}
	updateOnConflictValues := extractFilledFieldsFromModel(
		reflect.ValueOf(*logCount),
		reflect.TypeOf(*logCount),
	)

	// Upsert
	db = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "type"}}, // NOTE set to primary keys for table
		DoUpdates: clause.Assignments(updateOnConflictValues),
	}).Create(logCount)

	return db.Error
}

// StartLogCountLoader starts loader
func StartLogCountLoader() {
	go func() {
		postgresLoaderChan := GetLogCountModel().LoaderChannel

		for {
			// Read log
			newLogCount := <-postgresLoaderChan

			//////////////////////////
			// Get count from redis //
			//////////////////////////
			countKey := "log_count_" + newLogCount.Type

			count, err := redis.GetRedisClient().GetCount(countKey)
			if err != nil {
				zap.S().Fatal(
					"Loader=Log,",
					" Transaction Hash =", newLogCount.TransactionHash,
					" Log Index =", newLogCount.LogIndex,
					" Type=", newLogCount.Type,
					" - Error: ", err.Error())
			}

			// No count set yet
			// Get from database
			if count == -1 {
				curLogCount, err := GetLogCountModel().SelectOne(newLogCount.Type)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					count = 0
				} else if err != nil {
					zap.S().Fatal(
						"Loader=Log,",
						" Transaction Hash =", newLogCount.TransactionHash,
						" Log Index =", newLogCount.LogIndex,
						" Type=", newLogCount.Type,
						" - Error: ", err.Error())
				} else {
					count = int64(curLogCount.Count)
				}

				// Set count
				err = redis.GetRedisClient().SetCount(countKey, int64(count))
				if err != nil {
					// Redis error
					zap.S().Fatal(
						"Loader=Log,",
						" Transaction Hash =", newLogCount.TransactionHash,
						" Log Index =", newLogCount.LogIndex,
						" Type=", newLogCount.Type,
						" - Error: ", err.Error())
				}
			}

			//////////////////////
			// Load to postgres //
			//////////////////////

			// Add log to indexed
			if newLogCount.Type == "log" {
				newLogCountIndex := &models.LogCountIndex{
					TransactionHash: newLogCount.TransactionHash,
					LogIndex:        newLogCount.LogIndex,
				}
				err = GetLogCountIndexModel().Insert(newLogCountIndex)
				if err != nil {
					// Record already exists, continue
					continue
				}
			}

			// Increment records
			count, err = redis.GetRedisClient().IncCount(countKey)
			if err != nil {
				// Redis error
				zap.S().Fatal(
					"Loader=Log,",
					" Transaction Hash =", newLogCount.TransactionHash,
					" Log Index =", newLogCount.LogIndex,
					" Type=", newLogCount.Type,
					" - Error: ", err.Error())
			}
			newLogCount.Count = uint64(count)

			err = GetLogCountModel().UpsertOne(newLogCount)
			zap.S().Debug(
				"Loader=Log,",
				" Transaction Hash =", newLogCount.TransactionHash,
				" Log Index =", newLogCount.LogIndex,
				" Type=", newLogCount.Type,
				" - Upsert")
			if err != nil {
				// Postgres error
				zap.S().Fatal(
					"Loader=Log,",
					" Transaction Hash =", newLogCount.TransactionHash,
					" Log Index =", newLogCount.LogIndex,
					" Type=", newLogCount.Type,
					" - Error: ", err.Error())
			}
		}
	}()
}
