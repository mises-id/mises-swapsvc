package swap

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/mises-id/mises-swapsvc/app/models"
)

type (
	rpcError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	rpcBlockResponse struct {
		Result *models.BlockData `json:"result"`
		Error  *rpcError         `json:"error"`
	}
)

func getBlockByNumber(ctx context.Context, rpcAddress string, number string) (*models.BlockData, error) {
	api := rpcAddress
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	client.Timeout = time.Second * 3
	requestParam := make(map[string]interface{}, 1)
	requestParam["jsonrpc"] = "2.0"
	requestParam["id"] = "1"
	requestParam["method"] = "eth_getBlockByNumber"
	requestParam["params"] = []interface{}{number, true}
	jsonBytes, _ := json.Marshal(requestParam)
	req, err := http.NewRequest("POST", api, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	out := &rpcBlockResponse{}
	json.Unmarshal(body, out)
	if out.Error != nil {
		return nil, errors.New(out.Error.Message)
	}
	return out.Result, err
}

// getDecodedTransactionByHash
func getDecodedTransactionByHash(ctx context.Context, chain, hash string) (*models.TransactionDecodedReceipt, error) {
	api := fmt.Sprintf("https://deep-index.moralis.io/api/v2/transaction/%s/verbose?chain=%s", hash, chain)
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	client.Timeout = time.Second * 3
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("X-API-Key", "yGzL3EZY6u3MOBnu3im5N2JkJyY8v3x2NVhZk9GxHoQKOXdfYLu0PZj2JFP4Ygo2")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	out := &models.TransactionDecodedReceipt{}
	json.Unmarshal(body, out)
	if out.Message != "" {
		return nil, errors.New(out.Message)
	}
	return out, nil
}

func setProxy() func(*http.Request) (*url.URL, error) {
	return func(_ *http.Request) (*url.URL, error) {
		//return nil, nil
		return url.Parse("http://127.0.0.1:8889")
	}
}
