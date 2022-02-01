package crud

import (
	"reflect"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/geometry-labs/icon-logs/models"
)

// LogMissingModel - type for logMissing table model
type LogMissingModel struct {
	db            *gorm.DB
	model         *models.LogMissing
	modelORM      *models.LogMissingORM
	LoaderChannel chan *models.LogMissing
}

var logMissingModel *LogMissingModel
var logMissingModelOnce sync.Once

// GetLogMissingModel - create and/or return the logMissings table model
func GetLogMissingModel() *LogMissingModel {
	logMissingModelOnce.Do(func() {
		dbConn := getPostgresConn()
		if dbConn == nil {
			zap.S().Fatal("Cannot connect to postgres database")
		}

		logMissingModel = &LogMissingModel{
			db:            dbConn,
			model:         &models.LogMissing{},
			LoaderChannel: make(chan *models.LogMissing, 1),
		}

		err := logMissingModel.Migrate()
		if err != nil {
			zap.S().Fatal("LogMissingModel: Unable migrate postgres table: ", err.Error())
		}

		StartLogMissingLoader()
	})

	return logMissingModel
}

// Migrate - migrate logMissings table
func (m *LogMissingModel) Migrate() error {
	// Only using LogMissingRawORM (ORM version of the proto generated struct) to create the TABLE
	err := m.db.AutoMigrate(m.modelORM) // Migration and Index creation
	return err
}

func (m *LogMissingModel) FindMissing() error {
	db := m.db

	db.Exec(`
		CREATE TABLE log_missing_by_block_number AS
		SELECT
			transaction_hash, block_number, max_logs, num_logs
		FROM (
			SELECT
				transaction_hash,
		  	count(transaction_hash) as num_logs,
		    max(max_log_index) as max_logs,
				max(block_number) as block_number
		  FROM
				logs
		  GROUP BY
				transaction_hash
		) AS ml WHERE num_logs != max_logs;
	`,
	)

	return db.Error
}

func (m *LogMissingModel) UpsertOne(
	logMissing *models.LogMissing,
) error {
	db := m.db

	// map[string]interface{}
	updateOnConflictValues := extractFilledFieldsFromModel(
		reflect.ValueOf(*logMissing),
		reflect.TypeOf(*logMissing),
	)

	// Upsert
	db = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "transaction_hash"}}, // NOTE set to primary keys for table
		DoUpdates: clause.Assignments(updateOnConflictValues),
	}).Create(logMissing)

	return db.Error
}

// StartLogMissingLoader starts loader
func StartLogMissingLoader() {
	go func() {

		for {
			// Read transaction
			newLogMissing := <-GetLogMissingModel().LoaderChannel

			//////////////////////
			// Load to postgres //
			//////////////////////
			err := GetLogMissingModel().UpsertOne(newLogMissing)
			zap.S().Debug("Loader=LogMissing, TransactionHash=", newLogMissing.TransactionHash, " - Upsert")
			if err != nil {
				// Postgres error
				zap.S().Fatal("Loader=LogMissing, TransactionHash=", newLogMissing.TransactionHash, " - Error: ", err.Error())
			}
		}
	}()
}
