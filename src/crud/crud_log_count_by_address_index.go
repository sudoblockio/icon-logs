package crud

import (
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/geometry-labs/icon-logs/models"
)

// LogCountByAddressIndexModel - type for address table model
type LogCountByAddressIndexModel struct {
	db            *gorm.DB
	model         *models.LogCountByAddressIndex
	modelORM      *models.LogCountByAddressIndexORM
	LoaderChannel chan *models.LogCountByAddressIndex
}

var logCountByAddressIndexModel *LogCountByAddressIndexModel
var logCountByAddressIndexModelOnce sync.Once

// GetAddressModel - create and/or return the addresss table model
func GetLogCountByAddressIndexModel() *LogCountByAddressIndexModel {
	logCountByAddressIndexModelOnce.Do(func() {
		dbConn := getPostgresConn()
		if dbConn == nil {
			zap.S().Fatal("Cannot connect to postgres database")
		}

		logCountByAddressIndexModel = &LogCountByAddressIndexModel{
			db:            dbConn,
			model:         &models.LogCountByAddressIndex{},
			LoaderChannel: make(chan *models.LogCountByAddressIndex, 1),
		}

		err := logCountByAddressIndexModel.Migrate()
		if err != nil {
			zap.S().Fatal("LogCountByAddressIndexModel: Unable migrate postgres table: ", err.Error())
		}
	})

	return logCountByAddressIndexModel
}

// Migrate - migrate logCountByAddressIndexs table
func (m *LogCountByAddressIndexModel) Migrate() error {
	// Only using LogCountByAddressIndexRawORM (ORM version of the proto generated struct) to create the TABLE
	err := m.db.AutoMigrate(m.modelORM) // Migration and Index creation
	return err
}

// Insert - Insert logCountByIndex into table
func (m *LogCountByAddressIndexModel) Insert(logCountByAddressIndex *models.LogCountByAddressIndex) error {
	db := m.db

	// Set table
	db = db.Model(&models.LogCountByAddressIndex{})

	db = db.Create(logCountByAddressIndex)

	return db.Error
}
