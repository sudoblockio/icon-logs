package fixtures

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/geometry-labs/icon-logs/models"
)

const (
	logRawFixturesPath = "logs_raw.json"
)

// Fixtures - slice of Fixture
type Fixtures []Fixture

// Fixture - loaded from fixture file
type Fixture map[string]interface{}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// LoadLogFixtures - load log fixtures from disk
func LoadLogFixtures() []*models.Log {
	logs := make([]*models.Log, 0)

	fixtures, err := loadFixtures(logRawFixturesPath)
	check(err)

	for _, fixture := range fixtures {
		logs = append(logs, parseFixtureToLog(fixture))
	}

	return logs
}

func loadFixtures(file string) (Fixtures, error) {
	var fs Fixtures

	dat, err := ioutil.ReadFile(getFixtureDir() + file)
	check(err)
	err = json.Unmarshal(dat, &fs)

	return fs, err
}

func getFixtureDir() string {

	callDir, _ := os.Getwd()
	callDirSplit := strings.Split(callDir, "/")

	for i := len(callDirSplit) - 1; i >= 0; i-- {
		if callDirSplit[i] != "src" {
			callDirSplit = callDirSplit[:len(callDirSplit)-1]
		} else {
			break
		}
	}

	callDirSplit = append(callDirSplit, "fixtures")
	fixtureDir := strings.Join(callDirSplit, "/")
	fixtureDir = fixtureDir + "/"
	zap.S().Info(fixtureDir)

	return fixtureDir
}

func parseFixtureToLog(m map[string]interface{}) *models.Log {

  return &models.Log {
    Type: m["type"].(string),
    LogIndex: m["log_index"].(uint64),
    TransactionHash: m["transaction_hash"].(string),
    TransactionIndex: m["transaction_index"].(uint32),
    Address: m["address"].(string)
    Data: m["data"].(string)
    Indexed: m["indexed"].(string)
    BlockNumber: m["block_number"].(uint64)
    BlockTimestamp: m["block_timestamp"].(uint64)
    BlockHash: m["block_hash"].(string),
    ItemId: m["item_id"].(string),
    ItemTimestamp: m["item_timestamp"].(string),
  }
}
