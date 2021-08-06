package crud

import (
	"testing"
	"time"
  "context"

	"github.com/stretchr/testify/assert"

	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/fixtures"
	"github.com/geometry-labs/icon-logs/logging"
)

func init() {
	// Read env
	// Defaults should work
	config.ReadEnvironment()

	// Set up logging
	logging.Init()
}

func TestGetTransactionModel(t *testing.T) {
	assert := assert.New(t)

	logModel := GetTransactionModel()
	assert.NotEqual(nil, logModel)
}

func TestTransactionModelInsert(t *testing.T) {
	assert := assert.New(t)

	logModel := GetTransactionModel()
	assert.NotEqual(nil, logModel)

	// Load fixtures
	logFixtures := fixtures.LoadTransactionFixtures()

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

	for _, tx := range logFixtures {

		insertErr := logModel.Insert(ctx, tx)
		assert.Equal(nil, insertErr)
	}
}

func TestTransactionModelSelect(t *testing.T) {
	assert := assert.New(t)

	logModel := GetTransactionModel()
	assert.NotEqual(nil, logModel)

	// Load fixtures
	logFixtures := fixtures.LoadTransactionFixtures()

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

	for _, tx := range logFixtures {

		insertErr := logModel.Insert(ctx, tx)
		assert.Equal(nil, insertErr)
	}

	// Select all logs
	logs, err := logModel.Select(ctx, int64(len(logFixtures)), 0, "")
	assert.Equal(len(logFixtures), len(logs))
  assert.Equal(nil, err)

	// Test limit
	logs, err = logModel.Select(ctx, 1, 0, "")
	assert.Equal(1, len(logs))
  assert.Equal(nil, err)

	// Test skip
	logs, err = logModel.Select(ctx, 1, 1, "")
	assert.Equal(1, len(logs))
  assert.Equal(nil, err)
}

func TestTransactionModelLoader(t *testing.T) {
	assert := assert.New(t)

	logModel := GetTransactionModel()
	assert.NotEqual(nil, logModel)

	// Load fixtures
	logFixtures := fixtures.LoadTransactionFixtures()

	// Start loader
	go StartTransactionLoader()

	// Write to loader channel
	go func() {
		for _, fixture := range logFixtures {
			logModel.WriteChan <- fixture
		}
	}()

	// Wait for inserts
	time.Sleep(5)

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

	// Select all logs
	logs, err := logModel.Select(ctx, int64(len(logFixtures)), 0, "", "")
	assert.Equal(len(logFixtures), len(logs))
  assert.Equal(nil, err)
}
