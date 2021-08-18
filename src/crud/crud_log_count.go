package crud

import (
	"strings"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/geometry-labs/icon-logs/models"
)

// LogCountModel - type for log table model
type LogCountModel struct {
	db        *gorm.DB
	model     *models.LogCount
	modelORM  *models.LogCountORM
	WriteChan chan *models.LogCount
}

var logCountModel *LogCountModel
var logCountModelOnce sync.Once

// GetLogModel - create and/or return the logs table model
func GetLogCountModel() *LogCountModel {
	logCountModelOnce.Do(func() {
		dbConn := getPostgresConn()
		if dbConn == nil {
			zap.S().Fatal("Cannot connect to postgres database")
		}

		logCountModel = &LogCountModel{
			db:        dbConn,
			model:     &models.LogCount{},
			WriteChan: make(chan *models.LogCount, 1),
		}

		err := logCountModel.Migrate()
		if err != nil {
			zap.S().Fatal("LogCountModel: Unable migrate postgres table: ", err.Error())
		}
	})

	return logCountModel
}

// Migrate - migrate logCounts table
func (m *LogCountModel) Migrate() error {
	// Only using LogCountRawORM (ORM version of the proto generated struct) to create the TABLE
	err := m.db.AutoMigrate(m.modelORM) // Migration and Index creation
	return err
}

// Insert - Insert logCount into table
func (m *LogCountModel) Insert(logCount *models.LogCount) error {

	err := backoff.Retry(func() error {
		query := m.db.Create(logCount)

		if query.Error != nil && !strings.Contains(query.Error.Error(), "duplicate key value violates unique constraint") {
			zap.S().Warn("POSTGRES Insert Error : ", query.Error.Error())
			return query.Error
		}

		return nil
	}, backoff.NewExponentialBackOff())

	return err
}

// Update - Update logCount
func (m *LogCountModel) Update(logCount *models.LogCount) error {

	err := backoff.Retry(func() error {
		query := m.db.Model(&models.LogCount{}).Where("id = ?", logCount.Id).Update("count", logCount.Count)

		if query.Error != nil && !strings.Contains(query.Error.Error(), "duplicate key value violates unique constraint") {
			zap.S().Warn("POSTGRES Insert Error : ", query.Error.Error())
			return query.Error
		}

		return nil
	}, backoff.NewExponentialBackOff())

	return err
}

// Select - select from logCounts table
func (m *LogCountModel) Select() (models.LogCount, error) {
	db := m.db

	logCount := models.LogCount{}
	db = db.First(&logCount)

	return logCount, db.Error
}

// Delete - delete from logCounts table
func (m *LogCountModel) Delete(logCount models.LogCount) error {
	db := m.db

	db = db.Delete(&logCount)

	return db.Error
}

// StartLogCountLoader starts loader
func StartLogCountLoader() {
	go func() {
		var logCount *models.LogCount
		postgresLoaderChan := GetLogCountModel().WriteChan

		for {
			// Read logCount
			logCount = <-postgresLoaderChan

			// Load logCount to database
			curCount, err := GetLogCountModel().Select()
			if err == nil {
				logCount.Count = logCount.Count + curCount.Count
				GetLogCountModel().Update(logCount)
			} else {
				GetLogCountModel().Insert(logCount)
			}

		}
	}()
}
