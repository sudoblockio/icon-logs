package fixtures

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geometry-labs/icon-logs/config"
	"github.com/geometry-labs/icon-logs/logging"
)

func init() {
	// Read env
	// Defaults should work
	config.ReadEnvironment()

	// Set up logging
	logging.Init()
}

func TestLoadTransactionFixtures(t *testing.T) {
	assert := assert.New(t)

	transactionFixtures := LoadTransactionFixtures()

	assert.NotEqual(0, len(transactionFixtures))
}
