package swap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mises-id/mises-swapsvc/lib/codes"
)

type (
	SwapQuoteInfo struct {
		Aggregator       *Aggregator `json:"aggregator"`       //聚合器信息
		FromTokenAddress string      `json:"fromTokenAddress"` //swap发起询价的token地址
		ToTokenAddress   string      `json:"toTokenAddress"`   //swap目标token
		FromTokenAmount  string      `json:"fromTokenAmount"`  //swap发起询价token的数目
		ToTokenAmount    string      `json:"toTokenAmount"`    //swap目标token的数目
		Fee              float32     `json:"fee"`              //收取的佣金百分比，FromTokenAmount令牌数量的这个百分比将被发送到 referrerAddress，其余部分将用作交换的输入
		EstimateGasFee   string      `json:"estimateGasFee"`
		Error            string      `json:"error"`
		FetchTime        int64       `json:"fetch_time"` //
	}
	oneInchQuoteResponse struct {
		StatusCode     string     `json:"statusCode"`
		Error          string     `json:"error"`
		Description    string     `json:"description"`
		ToTokenAmount  string     `json:"toTokenAmount"`
		EstimateGasFee int64      `json:"estimatedGas"`
		Tx             *oneInchTx `json:"tx"`
	}
)

func swapQuote(ctx context.Context, in *SwapQuoteInput) ([]*SwapQuoteInfo, error) {
	//preload input parameters
	if err := preloadSwapQuoteInput(in); err != nil {
		return nil, err
	}
	//find providers
	//1inch
	quotes := make([]*SwapQuoteInfo, 0)
	st := time.Now()
	quoteOneInch := swapQuoteByOneInchProvider(ctx, in)
	quoteOneInch.FetchTime = int64(time.Since(st).Milliseconds())
	quotes = append(quotes, quoteOneInch)
	return quotes, nil
}

// ----------------------------------------------------------------
func preloadSwapQuoteInput(in *SwapQuoteInput) error {
	if in.ChainID == 0 {
		return codes.ErrInvalidArgument.New("Invaild chainID")
	}
	if in.FromTokenAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild fromTokenAddress")
	}
	if in.Amount == "" || in.Amount == "0" {
		return codes.ErrInvalidArgument.New("Invaild amount")
	}
	if in.ToTokenAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild toTokenAddress")
	}
	return nil
}

func swapQuoteByOneInchProvider(ctx context.Context, in *SwapQuoteInput) *SwapQuoteInfo {
	//check if OneInch is supported
	resp := &SwapQuoteInfo{
		FromTokenAddress: in.FromTokenAddress,
		ToTokenAddress:   in.ToTokenAddress,
		Fee:              swapFee,
		FromTokenAmount:  in.Amount,
	}
	resp.Aggregator = getAggregatorByProviderKey(ctx, oneInchProvider)
	data, err := apiSwapQuoteByOneInchProvider(ctx, in)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	resp.EstimateGasFee = fmt.Sprintf("%d", data.EstimateGasFee)
	resp.Error = data.Description
	resp.ToTokenAmount = data.ToTokenAmount
	return resp
}

func apiSwapQuoteByOneInchProvider(ctx context.Context, in *SwapQuoteInput) (*oneInchQuoteResponse, error) {
	api := fmt.Sprintf("%s/%d/quote?fromTokenAddress=%s&toTokenAddress=%s&amount=%s&fee=%.3f", oneIncheProviderAPIBaseURL, in.ChainID, in.FromTokenAddress, in.ToTokenAddress, in.Amount, swapFee)
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
		if strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
			return nil, errors.New("Request oneInch API timeout")
		}
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//return nil, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	out := &oneInchQuoteResponse{}
	json.Unmarshal(body, out)
	return out, nil
}
