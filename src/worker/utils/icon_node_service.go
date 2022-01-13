package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/geometry-labs/icon-logs/config"
)

func IconNodeServiceGetBlockTransactionHashes(height int) (*[]string, error) {

	// Request icon contract
	url := config.Config.IconNodeServiceURL
	method := "POST"
	payload := fmt.Sprintf(`{
    "jsonrpc": "2.0",
    "method": "icx_getBlockByHeight",
    "id": 1,
    "params": {
        "height": "0x%x"
    }
	}`, height)

	// Create http client
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// Execute request
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read body
	bodyString, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Check status code
	if res.StatusCode != 200 {
		return nil, errors.New(
			"StatusCode=" + strconv.Itoa(res.StatusCode) +
				",Request=" + payload +
				",Response=" + string(bodyString),
		)
	}

	// Parse body
	body := map[string]interface{}{}
	err = json.Unmarshal(bodyString, &body)
	if err != nil {
		return nil, err
	}

	// Extract result
	result, ok := body["result"].(map[string]interface{})
	if ok == false {
		return nil, errors.New("Cannot read result")
	}

	// Extract transactions
	transactions, ok := result["confirmed_transaction_list"].([]interface{})
	if ok == false {
		return nil, errors.New("Cannot read confirmed_transaction_list")
	}

	// Extract transaciton hashes
	transactionHashes := []string{}
	for _, t := range transactions {
		tx, ok := t.(map[string]interface{})
		if ok == false {
			return nil, errors.New("1 Cannot read transaction hash from block #" + fmt.Sprintf("%x", height))
		}

		// V1
		hash, ok := tx["tx_hash"].(string)
		if ok == true {
			transactionHashes = append(transactionHashes, "0x"+hash)
			continue
		}

		// V3
		hash, ok = tx["txHash"].(string)
		if ok == true {
			transactionHashes = append(transactionHashes, hash)
			continue
		}

		if ok == false {
			return nil, errors.New("2 Cannot read transaction hash from block #" + fmt.Sprintf("%x", height))
		}
	}

	return &transactionHashes, nil
}

func IconNodeServiceGetTransactionLogCount(transaction_hash string) (int, error) {

	// Request icon contract
	url := config.Config.IconNodeServiceURL
	method := "POST"
	payload := fmt.Sprintf(`{
    "jsonrpc": "2.0",
    "method": "icx_getTransactionResult",
    "id": 1234,
    "params": {
        "txHash": "%s"
    }	
	}`, transaction_hash)

	// Create http client
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return 0, err
	}

	// Execute request
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// Read body
	bodyString, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	// Check status code
	if res.StatusCode != 200 {
		return 0, errors.New(
			"StatusCode=" + strconv.Itoa(res.StatusCode) +
				",Request=" + payload +
				",Response=" + string(bodyString),
		)
	}

	// Parse body
	body := map[string]interface{}{}
	err = json.Unmarshal(bodyString, &body)
	if err != nil {
		return 0, err
	}

	// Extract result
	result, ok := body["result"].(map[string]interface{})
	if ok == false {
		return 0, errors.New("Cannot read result")
	}

	// Extract logs
	logs, ok := result["eventLogs"].([]interface{})
	if ok == false {
		return 0, errors.New("Cannot read event logs")
	}

	return len(logs), nil
}
