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

// LogCountByAddressModel - type for address table model
type LogCountByAddressModel struct {
	db            *gorm.DB
	model         *models.LogCountByAddress
	modelORM      *models.LogCountByAddressORM
	LoaderChannel chan *models.LogCountByAddress
}

var logCountByAddressModel *LogCountByAddressModel
var logCountByAddressModelOnce sync.Once

// GetAddressModel - create and/or return the addresss table model
func GetLogCountByAddressModel() *LogCountByAddressModel {
	logCountByAddressModelOnce.Do(func() {
		dbConn := getPostgresConn()
		if dbConn == nil {
			zap.S().Fatal("Cannot connect to postgres database")
		}

		logCountByAddressModel = &LogCountByAddressModel{
			db:            dbConn,
			model:         &models.LogCountByAddress{},
			LoaderChannel: make(chan *models.LogCountByAddress, 1),
		}

		err := logCountByAddressModel.Migrate()
		if err != nil {
			zap.S().Fatal("LogCountByAddressModel: Unable migrate postgres table: ", err.Error())
		}

		StartLogCountByAddressLoader()
	})

	return logCountByAddressModel
}

// Migrate - migrate logCountByAddresss table
func (m *LogCountByAddressModel) Migrate() error {
	// Only using LogCountByAddressRawORM (ORM version of the proto generated struct) to create the TABLE
	err := m.db.AutoMigrate(m.modelORM) // Migration and Index creation
	return err
}

// Select - select from logCountByAddresss table
func (m *LogCountByAddressModel) SelectOne(address string) (*models.LogCountByAddress, error) {
	db := m.db

	// Set table
	db = db.Model(&models.LogCountByAddress{})

	// Address
	db = db.Where("address = ?", address)

	logCountByAddress := &models.LogCountByAddress{}
	db = db.First(logCountByAddress)

	return logCountByAddress, db.Error
}

// SelectMany - select from logCountByAddresss table
func (m *LogCountByAddressModel) SelectMany(limit int, skip int) (*[]models.LogCountByAddress, error) {
	db := m.db

	// Set table
	db = db.Model(&[]models.LogCountByAddress{})

	// Limit
	db = db.Limit(limit)

	// Skip
	if skip != 0 {
		db = db.Offset(skip)
	}

	logCountByAddresses := &[]models.LogCountByAddress{}
	db = db.Find(logCountByAddresses)

	return logCountByAddresses, db.Error
}

// Select - select from logCountByAddresss table
func (m *LogCountByAddressModel) SelectCount(address string) (uint64, error) {
	db := m.db

	// Set table
	db = db.Model(&models.LogCountByAddress{})

	// Address
	db = db.Where("address = ?", address)

	logCountByAddress := &models.LogCountByAddress{}
	db = db.First(logCountByAddress)

	count := uint64(0)
	if logCountByAddress != nil {
		count = logCountByAddress.Count
	}

	return count, db.Error
}

func (m *LogCountByAddressModel) UpsertOne(
	logCountByAddress *models.LogCountByAddress,
) error {
	db := m.db

	// map[string]interface{}
	updateOnConflictValues := extractFilledFieldsFromModel(
		reflect.ValueOf(*logCountByAddress),
		reflect.TypeOf(*logCountByAddress),
	)

	// Upsert
	db = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}}, // NOTE set to primary keys for table
		DoUpdates: clause.Assignments(updateOnConflictValues),
	}).Create(logCountByAddress)

	return db.Error
}

// StartLogCountByAddressLoader starts loader
func StartLogCountByAddressLoader() {
	go func() {
		postgresLoaderChan := GetLogCountByAddressModel().LoaderChannel

		for {
			// Read log
			newLogCountByAddress := <-postgresLoaderChan

			//////////////////////////
			// Get count from redis //
			//////////////////////////
			countKey := "icon_logs_log_count_by_address_" + newLogCountByAddress.Address

			count, err := redis.GetRedisClient().GetCount(countKey)
			if err != nil {
				zap.S().Fatal(
					"Loader=Log",
					" Transaction Hash=", newLogCountByAddress.TransactionHash,
					" Log Index=", newLogCountByAddress.LogIndex,
					" Address=", newLogCountByAddress.Address,
					" - Error: ", err.Error())
			}

			// No count set yet
			// Get from database
			if count == -1 {
				curLogCountByAddress, err := GetLogCountByAddressModel().SelectOne(newLogCountByAddress.Address)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					count = 0
				} else if err != nil {
					zap.S().Fatal(
						"Loader=Log",
						" Transaction Hash=", newLogCountByAddress.TransactionHash,
						" Log Index=", newLogCountByAddress.LogIndex,
						" Address=", newLogCountByAddress.Address,
						" - Error: ", err.Error())
				} else {
					count = int64(curLogCountByAddress.Count)
				}

				// Set count
				err = redis.GetRedisClient().SetCount(countKey, int64(count))
				if err != nil {
					// Redis error

					zap.S().Fatal(
						"Loader=Log",
						" Transaction Hash=", newLogCountByAddress.TransactionHash,
						" Log Index=", newLogCountByAddress.LogIndex,
						" Address=", newLogCountByAddress.Address,
						" - Error: ", err.Error())
				}
			}

			//////////////////////
			// Load to postgres //
			//////////////////////

			// Add log to indexed
			newLogCountByAddressIndex := &models.LogCountByAddressIndex{
				TransactionHash: newLogCountByAddress.TransactionHash,
				LogIndex:        newLogCountByAddress.LogIndex,
			}
			err = GetLogCountByAddressIndexModel().Insert(newLogCountByAddressIndex)
			if err != nil {
				// Record already exists, continue
				continue
			}

			// Increment records
			count, err = redis.GetRedisClient().IncCount(countKey)
			if err != nil {
				// Redis error

				zap.S().Fatal(
					"Loader=Log",
					" Transaction Hash=", newLogCountByAddress.TransactionHash,
					" Log Index=", newLogCountByAddress.LogIndex,
					" Address=", newLogCountByAddress.Address,
					" - Error: ", err.Error())
			}
			newLogCountByAddress.Count = uint64(count)

			err = GetLogCountByAddressModel().UpsertOne(newLogCountByAddress)
			zap.S().Debug("Loader=Log, Hash=", newLogCountByAddress.TransactionHash, " Address=", newLogCountByAddress.Address, " - Upserted")
			if err != nil {
				// Postgres error

				zap.S().Fatal(
					"Loader=Log",
					" Transaction Hash=", newLogCountByAddress.TransactionHash,
					" Log Index=", newLogCountByAddress.LogIndex,
					" Address=", newLogCountByAddress.Address,
					" - Error: ", err.Error())
			}
		}
	}()
}
