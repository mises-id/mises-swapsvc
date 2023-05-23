package swap_sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/mises-id/mises-swapsvc/app/models"
	"github.com/mises-id/mises-swapsvc/lib/utils"
)

type (
	GetTransactionsByAddressRequest struct {
		Address    string
		StartBlock uint64
		EndBlock   uint64
		Page       uint64
		Offset     uint64
		Sort       string
		ApiKey     string
	}
	GetTransactionsByAddressResponse struct {
		Status  string                `json:"status"`
		Message string                `json:"message"`
		Result  []*models.Transaction `json:"result"`
	}
	EthBlockNumberResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}
)

func apiGetTransactionsByAddress(ctx context.Context, endpoint string, apiKey string, params *GetTransactionsByAddressRequest) ([]*models.Transaction, error) {
	if params == nil {
		return nil, errors.New("Invalid request params")
	}
	api := fmt.Sprintf("%sapi?module=account&action=txlist&address=%s&startblock=%d&endblock=%d&page=%d&offset=%d&sort=%s&apikey=%s", endpoint, params.Address, params.StartBlock, params.EndBlock, params.Page, params.Offset, params.Sort, apiKey)
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	client.Timeout = time.Second * 3
	req, err := http.NewRequest("GET", api, nil)
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
	out := &GetTransactionsByAddressResponse{}
	json.Unmarshal(body, out)
	if out.Status == "0" {
		return nil, errors.New(out.Message)
	}
	return out.Result, err
}

func apiEthBlockNumber(ctx context.Context, endpoint, apiKey string) (uint64, error) {
	api := fmt.Sprintf("%sapi?module=proxy&action=eth_blockNumber&apikey=%s", endpoint, apiKey)
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	client.Timeout = time.Second * 5
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return 0, err
	}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	out := &EthBlockNumberResponse{}
	json.Unmarshal(body, out)
	if out.Status == "0" {
		return 0, errors.New(out.Message)
	}
	if out.Result == "" {
		return 0, errors.New("eth_blockNumber no result")
	}
	return utils.Hex2Dec(out.Result), nil
}

func setProxy() func(*http.Request) (*url.URL, error) {
	return func(_ *http.Request) (*url.URL, error) {
		return nil, nil
		return url.Parse("http://127.0.0.1:8889")
	}
}
