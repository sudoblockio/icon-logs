package crud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/logging"
	"github.com/geometry-labs/icon-logs/models"
)

func init() {
	// Read env
	// Defaults should work
	config.ReadEnvironment()

	// Set up logging
	logging.Init()
}

func TestGetLogCountModel(t *testing.T) {
	assert := assert.New(t)

	logCountModel := GetLogCountModel()
	assert.NotEqual(nil, logCountModel)
}

func TestLogCountModelInsert(t *testing.T) {
	assert := assert.New(t)

	logCountModel := GetLogCountModel()
	assert.NotEqual(nil, logCountModel)

	logCountFixture := &models.LogCount{
		Count: 1,
		Id:    10,
	}

	// Clear entry
	logCountModel.Delete(*logCountFixture)

	insertErr := logCountModel.Insert(logCountFixture)
	assert.Equal(nil, insertErr)
}

func TestLogCountModelSelect(t *testing.T) {
	assert := assert.New(t)

	logCountModel := GetLogCountModel()
	assert.NotEqual(nil, logCountModel)

	// Load fixture
	logCountFixture := &models.LogCount{
		Count: 1,
		Id:    10,
	}

	// Clear entry
	logCountModel.Delete(*logCountFixture)

	insertErr := logCountModel.Insert(logCountFixture)
	assert.Equal(nil, insertErr)

	// Select LogCount
	result, err := logCountModel.Select()
	assert.Equal(logCountFixture.Count, result.Count)
	assert.Equal(logCountFixture.Id, result.Id)
	assert.Equal(nil, err)
}

func TestLogCountModelUpdate(t *testing.T) {
	assert := assert.New(t)

	logCountModel := GetLogCountModel()
	assert.NotEqual(nil, logCountModel)

	// Load fixture
	logCountFixture := &models.LogCount{
		Count: 1,
		Id:    10,
	}

	// Clear entry
	logCountModel.Delete(*logCountFixture)

	insertErr := logCountModel.Insert(logCountFixture)
	assert.Equal(nil, insertErr)

	// Select LogCount
	result, err := logCountModel.Select()
	assert.Equal(logCountFixture.Count, result.Count)
	assert.Equal(logCountFixture.Id, result.Id)
	assert.Equal(nil, err)

	// Update LogCount
	logCountFixture = &models.LogCount{
		Count: 10,
		Id:    10,
	}
	insertErr = logCountModel.Update(logCountFixture)
	assert.Equal(nil, insertErr)

	// Select LogCount
	result, err = logCountModel.Select()
	assert.Equal(logCountFixture.Count, result.Count)
	assert.Equal(logCountFixture.Id, result.Id)
	assert.Equal(nil, err)
}

func TestLogCountModelLoader(t *testing.T) {
	assert := assert.New(t)

	logCountModel := GetLogCountModel()
	assert.NotEqual(nil, logCountModel)

	// Load fixture
	logCountFixture := &models.LogCount{
		Count: 1,
		Id:    10,
	}

	// Clear entry
	logCountModel.Delete(*logCountFixture)

	// Start loader
	StartLogCountLoader()

	// Write to loader channel
	go func() {
		for {
			logCountModel.WriteChan <- logCountFixture
			time.Sleep(1)
		}
	}()

	// Wait for inserts
	time.Sleep(5)
}
