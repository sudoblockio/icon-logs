package crud

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/geometry-labs/icon-logs/models"
)

func convertMapToBsonD(v map[string]interface{}) (*bson.D, error) {
  var doc *bson.D
  var err error

	data, err := bson.Marshal(v)
	if err != nil {
		return doc, err
	}

	err = bson.Unmarshal(data, &doc)
	return doc, err
}

func convertBsonMToLog(m bson.M) models.Log {

  // Data field may be null
  data, ok := m["data"].(string)
  if ok == false {
    data = ""
  }

  return models.Log {
    Type: m["type"].(string),
    LogIndex: m["logindex"].(uint64),
    TransactionHash: m["transactionhash"].(string),
    TransactionIndex: m["transactionindex"].(uint32),
    Address: m["address"].(string),
    Data: m["data"].(string),
    Indexed: m["indexed"].(string),
    BlockNumber: m["blocknumber"].(uint64),
    BlockTimestamp: m["blocktimestamp"].(uint64),
    BlockHash: m["blockhash"].(string),
    ItemId: m["itemid"].(string),
    ItemTimestamp: m["itemtimestamp"].(string),
  }
}
