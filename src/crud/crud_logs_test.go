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

func TestGetLogModel(t *testing.T) {
	assert := assert.New(t)

	logModel := GetLogModel()
	assert.NotEqual(nil, logModel)
}

func TestLogModelInsert(t *testing.T) {
	assert := assert.New(t)

	logModel := GetLogModel()
	assert.NotEqual(nil, logModel)

	// Load fixtures
	logFixtures := fixtures.LoadLogFixtures()

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

	for _, tx := range logFixtures {

		insertErr := logModel.Insert(ctx, tx)
		assert.Equal(nil, insertErr)
	}
}

func TestLogModelSelect(t *testing.T) {
	assert := assert.New(t)

	logModel := GetLogModel()
	assert.NotEqual(nil, logModel)

	// Load fixtures
	logFixtures := fixtures.LoadLogFixtures()

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

	for _, tx := range logFixtures {

		insertErr := logModel.Insert(ctx, tx)
		assert.Equal(nil, insertErr)
	}

	// Select all logs
	logs, err := logModel.Select(ctx, int64(len(logFixtures)), 0, "", "")
	assert.Equal(len(logFixtures), len(logs))
  assert.Equal(nil, err)

	// Test limit
	logs, err = logModel.Select(ctx, 1, 0, "", "")
	assert.Equal(1, len(logs))
  assert.Equal(nil, err)

	// Test skip
	logs, err = logModel.Select(ctx, 1, 1, "", "")
	assert.Equal(1, len(logs))
  assert.Equal(nil, err)

	// Test from
	logs, err = logModel.Select(ctx, 1, 0, "hx02e6bf5860b7d7744ec5050545d10d37c72ac2ef", "")
	assert.Equal(1, len(logs))
  assert.Equal(nil, err)

	// Test to
	logs, err = logModel.Select(ctx, 1, 0, "", "cx38fd2687b202caf4bd1bda55223578f39dbb6561")
	assert.Equal(1, len(logs))
  assert.Equal(nil, err)
}

func TestLogModelLoader(t *testing.T) {
	assert := assert.New(t)

	logModel := GetLogModel()
	assert.NotEqual(nil, logModel)

	// Load fixtures
	logFixtures := fixtures.LoadLogFixtures()

	// Start loader
	go StartLogLoader()

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
