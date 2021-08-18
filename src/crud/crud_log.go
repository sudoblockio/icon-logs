package crud

import (
	"strings"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/geometry-labs/icon-logs/models"
)

// LogModel - type for log table model
type LogModel struct {
	db        *gorm.DB
	model     *models.Log
	modelORM  *models.LogORM
	WriteChan chan *models.Log
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
			db:        dbConn,
			model:     &models.Log{},
			WriteChan: make(chan *models.Log, 1),
		}

		err := logModel.Migrate()
		if err != nil {
			zap.S().Fatal("LogModel: Unable migrate postgres table: ", err.Error())
		}
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

	err := backoff.Retry(func() error {
		query := m.db.Create(log)
		if query.Error != nil && !strings.Contains(query.Error.Error(), "duplicate key value violates unique constraint") {
			zap.S().Warn("POSTGRES Insert Error : ", query.Error.Error())
			return query.Error
		}

		return nil
	}, backoff.NewExponentialBackOff())

	return err
}

// Select - select from logs table
// Returns: models, total count (if filters), error (if present)
func (m *LogModel) Select(
	limit int,
	skip int,
	txHash string,
	scoreAddr string,
) ([]models.Log, int64, error) {
	db := m.db
	computeCount := false

	// Set table
	db = db.Model(&models.Log{})

	// Latest logs first
	db = db.Order("block_number desc")

	// Hash
	if txHash != "" {
		computeCount = true
		db = db.Where("transaction_hash = ?", txHash)
	}

	// Address
	if scoreAddr != "" {
		computeCount = true
		db = db.Where("address = ?", scoreAddr)
	}

	// Count, if needed
	count := int64(-1)
	if computeCount {
		db.Count(&count)
	}

	// Limit is required and defaulted to 1
	// Note: Count before setting limit
	db = db.Limit(limit)

	// Skip
	// Note: Count before setting skip
	if skip != 0 {
		db = db.Offset(skip)
	}

	logs := []models.Log{}
	db = db.Find(&logs)

	return logs, count, db.Error
}

// StartLogLoader starts loader
func StartLogLoader() {
	go func() {
		var log *models.Log
		postgresLoaderChan := GetLogModel().WriteChan

		for {
			// Read log
			log = <-postgresLoaderChan

			// Load log to database
			GetLogModel().Insert(log)

			zap.S().Debugf("Loader Log: Loaded in postgres table Logs, Block Number: %d", log.BlockNumber)
		}
	}()
}
