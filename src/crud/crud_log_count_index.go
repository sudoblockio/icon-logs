package crud

import (
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/geometry-labs/icon-logs/models"
)

// LogCountIndexModel - type for address table model
type LogCountIndexModel struct {
	db            *gorm.DB
	model         *models.LogCountIndex
	modelORM      *models.LogCountIndexORM
	LoaderChannel chan *models.LogCountIndex
}

var logCountIndexModel *LogCountIndexModel
var logCountIndexModelOnce sync.Once

// GetAddressModel - create and/or return the addresss table model
func GetLogCountIndexModel() *LogCountIndexModel {
	logCountIndexModelOnce.Do(func() {
		dbConn := getPostgresConn()
		if dbConn == nil {
			zap.S().Fatal("Cannot connect to postgres database")
		}

		logCountIndexModel = &LogCountIndexModel{
			db:            dbConn,
			model:         &models.LogCountIndex{},
			LoaderChannel: make(chan *models.LogCountIndex, 1),
		}

		err := logCountIndexModel.Migrate()
		if err != nil {
			zap.S().Fatal("LogCountIndexModel: Unable migrate postgres table: ", err.Error())
		}
	})

	return logCountIndexModel
}

// Migrate - migrate logCountIndexs table
func (m *LogCountIndexModel) Migrate() error {
	// Only using LogCountIndexRawORM (ORM version of the proto generated struct) to create the TABLE
	err := m.db.AutoMigrate(m.modelORM) // Migration and Index creation
	return err
}

// Count - count all entries in log_count_indices table
// NOTE this function will take a long time
func (m *LogCountIndexModel) Count() (int64, error) {
	db := m.db

	// Set table
	db = db.Model(&models.LogCountIndex{})

	// Count
	var count int64
	db = db.Count(&count)

	return count, db.Error
}

// Insert - Insert logCountByIndex into table
func (m *LogCountIndexModel) Insert(logCountIndex *models.LogCountIndex) error {
	db := m.db

	// Set table
	db = db.Model(&models.LogCountIndex{})

	db = db.Create(logCountIndex)

	return db.Error
}
