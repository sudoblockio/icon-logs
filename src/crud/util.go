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
  return models.Log {
    Type: m["type"].(string),
    LogIndex: uint64(m["logindex"].(int64)),
    TransactionHash: m["transactionhash"].(string),
    TransactionIndex: uint32(m["transactionindex"].(int64)),
    Address: m["address"].(string),
    Data: m["data"].(string),
    Indexed: m["indexed"].(string),
    BlockNumber: uint64(m["blocknumber"].(int64)),
    BlockTimestamp: uint64(m["blocktimestamp"].(int64)),
    BlockHash: m["blockhash"].(string),
    ItemId: m["itemid"].(string),
    ItemTimestamp: m["itemtimestamp"].(string),
  }
}
