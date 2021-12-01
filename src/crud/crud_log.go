package crud

import (
	"reflect"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/geometry-labs/icon-logs/models"
)

// LogModel - type for log table model
type LogModel struct {
	db            *gorm.DB
	model         *models.Log
	modelORM      *models.LogORM
	LoaderChannel chan *models.Log
}

var logModel *LogModel
var logModelOnce sync.Once

// GetLogModel - create and/or return the logs table model
func GetLogModel() *LogModel {
	logModelOnce.Do(func() {
		dbConn := getPostgresConn()
		if dbConn == nil {
			zap.S().Fatal("Cannot connect to postgres database")
		}

		logModel = &LogModel{
			db:            dbConn,
			model:         &models.Log{},
			LoaderChannel: make(chan *models.Log, 1),
		}

		err := logModel.Migrate()
		if err != nil {
			zap.S().Fatal("LogModel: Unable migrate postgres table: ", err.Error())
		}

		StartLogLoader()
	})

	return logModel
}

// Migrate - migrate logs table
func (m *LogModel) Migrate() error {
	// Only using LogRawORM (ORM version of the proto generated struct) to create the TABLE
	err := m.db.AutoMigrate(m.modelORM) // Migration and Index creation
	return err
}

// Insert - Insert log into table
func (m *LogModel) Insert(log *models.Log) error {
	db := m.db

	// Set table
	db = db.Model(&models.Log{})

	db = db.Create(log)

	return db.Error
}

// Select - select from logs table
// Returns: models, error (if present)
func (m *LogModel) SelectMany(
	limit int,
	skip int,
	blockNumber uint32,
	transactionHash string,
	scoreAddress string,
	method string,
) (*[]models.Log, error) {
	db := m.db

	// Set table
	db = db.Model(&models.Log{})

	// Latest logs first
	db = db.Order("block_number desc")

	// Number
	if blockNumber != 0 {
		db = db.Where("block_number = ?", blockNumber)
	}

	// Hash
	if transactionHash != "" {
		db = db.Where("transaction_hash = ?", transactionHash)
	}

	// Address
	if scoreAddress != "" {
		db = db.Where("address = ?", scoreAddress)
	}

	// Method
	if method != "" {
		db = db.Where("method = ?", method)
	}

	// Limit is required and defaulted to 1
	db = db.Limit(limit)

	// Skip
	if skip != 0 {
		db = db.Offset(skip)
	}

	logs := &[]models.Log{}
	db = db.Find(logs)

	return logs, db.Error
}

func (m *LogModel) SelectOne(
	transactionHash string,
	logIndex uint64,
) (*models.Log, error) {
	db := m.db

	// Set table
	db = db.Model(&models.Log{})

	db = db.Where("transaction_hash = ?", transactionHash)

	db = db.Where("log_index = ?", logIndex)

	log := &models.Log{}
	db = db.First(log)

	return log, db.Error
}

// CountByAddress - Count logs by address
// NOTE this function may take very long for some addresses
func (m *LogModel) CountByAddress(address string) (int64, error) {
	db := m.db

	// Set table
	db = db.Model(&models.Log{})

	// Address
	db = db.Where("address = ?", address)

	// Count
	var count int64
	db = db.Count(&count)

	return count, db.Error
}

// UpdateOne - select from logs table
func (m *LogModel) UpdateOne(
	log *models.Log,
) error {
	db := m.db

	// Set table
	db = db.Model(&models.Log{})

	// Transaction Hash
	db = db.Where("transaction_hash = ?", log.TransactionHash)

	// Log Index
	db = db.Where("log_index = ?", log.LogIndex)

	db = db.Save(log)

	return db.Error
}

func (m *LogModel) UpsertOne(
	log *models.Log,
) error {
	db := m.db

	// map[string]interface{}
	updateOnConflictValues := extractFilledFieldsFromModel(
		reflect.ValueOf(*log),
		reflect.TypeOf(*log),
	)

	// Upsert
	db = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "transaction_hash"}, {Name: "log_index"}}, // NOTE set to primary keys for table
		DoUpdates: clause.Assignments(updateOnConflictValues),
	}).Create(log)

	return db.Error
}

// StartLogLoader starts loader
func StartLogLoader() {
	go func() {

		for {
			// Read transaction
			newLog := <-GetLogModel().LoaderChannel

			//////////////////////
			// Load to postgres //
			//////////////////////
			err := GetLogModel().UpsertOne(newLog)
			zap.S().Debug("Loader=Log, TransactionHash=", newLog.TransactionHash, " LogIndex=", newLog.LogIndex, " - Upsert")
			if err != nil {
				// Postgres error
				zap.S().Fatal("Loader=Log, TransactionHash=", newLog.TransactionHash, " LogIndex=", newLog.LogIndex, " - Error: ", err.Error())
			}
		}
	}()
}
