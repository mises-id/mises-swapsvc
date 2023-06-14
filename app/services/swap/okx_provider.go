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
	OkxProviderConfig struct {
		swapFee                  float32
		minSlippage, maxSlippage float32
		slippageDecimals         float32
	}

	OkxProvider struct {
		*BaseProvider
		config *OkxProviderConfig
	}
	//api response
	OkxGetApproveAllowanceResponse struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data *struct {
			AllowanceAmount string `json:"allowanceAmount"`
		} `json:"data"`
	}
	OkxApproveTransactionResponse struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data *struct {
			Data               string `json:"data"`
			DexContractAddress string `json:"dexContractAddress"`
			GasLimit           string `json:"gasLimit"`
			GasPrice           string `json:"gasPrice"`
		} `json:"data"`
	}
	OkxSwapQuoteData struct {
		EstimateGasFee string `json:"estimateGasFee"`
		ToTokenAmount  string `json:"toTokenAmount"`
	}
	OkxSwapQuoteResponse struct {
		Code string              `json:"code"`
		Msg  string              `json:"msg"`
		Data []*OkxSwapQuoteData `json:"data"`
	}
	OkxTx struct {
		Data     string `json:"data"`
		From     string `json:"from"`
		To       string `json:"to"`
		GasPrice string `json:"gasPrice"`
		Gas      int64  `json:"gas"`
		Value    string `json:"value"`
	}
	OkxSwapTradeData struct {
		Tx            *OkxTx `json:"tx"`
		ToTokenAmount string `json:"toTokenAmount"`
	}
	OkxSwapTradeResponse struct {
		Code string              `json:"code"`
		Msg  string              `json:"msg"`
		Data []*OkxSwapTradeData `json:"data"`
	}
)

func NewOkxProvider() *OkxProvider {
	config := &OkxProviderConfig{
		swapFee:          0,
		minSlippage:      0.00001,
		maxSlippage:      1,
		slippageDecimals: 1,
	}

	resp := &OkxProvider{
		BaseProvider: NewBaseProvider(OkxProviderKey),
		config:       config,
	}
	return resp
}

// GetSwapApproveAllowance
func (p *OkxProvider) GetSwapApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error) {
	return p.getSwapApproveAllowance(ctx, in)
}

func (p *OkxProvider) getSwapApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error) {
	data, err := p.apiGetApproveAllowance(ctx, in)
	if err != nil {
		return nil, err
	}
	if data.Code != "0" || data.Data == nil {
		return nil, codes.ErrInvalidArgument.New(data.Msg)
	}
	res := &GetSwapApproveAllowanceOutput{
		Allowance: data.Data.AllowanceAmount,
	}
	return res, nil
}

// ApproveSwapTransaction
func (p *OkxProvider) ApproveSwapTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error) {
	return p.approveSwapTransaction(ctx, in)
}

func (p *OkxProvider) approveSwapTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error) {
	data, err := p.apiApproveTransaction(ctx, in)
	if err != nil {
		return nil, err
	}
	if data.Code != "0" || data.Data == nil {
		return nil, codes.ErrInvalidArgument.New(data.Msg)
	}
	res := &ApproveSwapTransactionOutput{
		Data:     data.Data.Data,
		GasPrice: data.Data.GasPrice,
		To:       data.Data.DexContractAddress,
	}
	return res, nil
}

// SwapQuotes
func (p *OkxProvider) SwapQuotes(ctx context.Context, in *SwapQuoteInput) []*SwapQuoteInfo {
	quotes := make([]*SwapQuoteInfo, 0)
	quote := p.swapQuote(ctx, in)
	if quote != nil {
		quotes = append(quotes, quote)
	}
	return quotes
}

// SwapQuote
func (p *OkxProvider) SwapQuote(ctx context.Context, in *SwapQuoteInput) *SwapQuoteInfo {
	return p.swapQuote(ctx, in)
}

func (p *OkxProvider) swapQuote(ctx context.Context, in *SwapQuoteInput) *SwapQuoteInfo {
	// check if OneInch is supported. chainId
	st := time.Now()
	resp := &SwapQuoteInfo{
		FromTokenAddress: in.FromTokenAddress,
		ToTokenAddress:   in.ToTokenAddress,
		Fee:              p.config.swapFee,
		FromTokenAmount:  in.Amount,
	}
	resp.Aggregator = p.aggregator(in.ChainID)
	contract := p.findContractByChainID(in.ChainID)
	if contract == nil {
		resp.Error = ErrUnspportChainID
		return resp
	}
	data, err := p.apiSwapQuote(ctx, in)
	resp.FetchTime = int64(time.Since(st).Milliseconds())
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	if data.Code != "0" || data.Data == nil || len(data.Data) == 0 {
		if data.Msg == "" {
			data.Msg = ErrorNotFoundQuote
		}
		resp.Error = data.Msg
		return resp
	}
	quote := data.Data[0]
	resp.EstimateGasFee = quote.EstimateGasFee
	resp.ToTokenAmount = quote.ToTokenAmount
	return resp
}

// SwapTrade
func (p *OkxProvider) SwapTrade(ctx context.Context, in *SwapTradeInput) (*SwapTradeInfo, error) {
	return p.swapTrade(ctx, in), nil
}

func (p *OkxProvider) swapTrade(ctx context.Context, in *SwapTradeInput) *SwapTradeInfo {
	if in == nil {
		return nil
	}
	//check if OneInch is supported
	st := time.Now()
	resp := &SwapTradeInfo{
		FromTokenAddress: in.FromTokenAddress,
		ToTokenAddress:   in.ToTokenAddress,
		Fee:              p.config.swapFee,
		FromTokenAmount:  in.Amount,
	}
	resp.Aggregator = p.aggregator(in.ChainID)
	contract := p.findContractByChainID(in.ChainID)
	if contract == nil {
		resp.Error = ErrUnspportChainID
		return resp
	}
	// handle DestReceiver
	in.DestReceiver = in.FromAddress
	data, err := p.apiSwapTrade(ctx, in)
	resp.FetchTime = int64(time.Since(st).Milliseconds())
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	if data.Code != "0" || data.Data == nil || len(data.Data) == 0 {
		if data.Msg == "" {
			data.Msg = ErrorNotFoundQuote
		}
		resp.Error = data.Msg
		return resp
	}
	trade := data.Data[0]
	resp.ToTokenAmount = trade.ToTokenAmount
	if trade.Tx != nil {
		resp.Trade = &Trade{
			From:     trade.Tx.From,
			To:       trade.Tx.To,
			Data:     trade.Tx.Data,
			GasPrice: trade.Tx.GasPrice,
			GasLimit: fmt.Sprintf("%d", trade.Tx.Gas),
			Value:    trade.Tx.Value,
		}
	}
	return resp
}

// ----------------------------------------------------------------
// api start
// ----------------------------------------------------------------
func (p *OkxProvider) apiGetApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*OkxGetApproveAllowanceResponse, error) {
	api := fmt.Sprintf("%s/dex/aggregator/get-allowance?chainId=%d&tokenContractAddress=%s&userWalletAddress=%s", OkxProviderAPIBaseURL, in.ChainID, in.TokenAddress, in.WalletAddress)
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
	out := &OkxGetApproveAllowanceResponse{}
	json.Unmarshal(body, out)
	return out, nil
}

func (p *OkxProvider) apiApproveTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*OkxApproveTransactionResponse, error) {
	api := fmt.Sprintf("%s/dex/aggregator/approve-transaction?chainId=%d&tokenContractAddress=%s&approveAmount=%s", OkxProviderAPIBaseURL, in.ChainID, in.TokenAddress, in.Amount)
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
	out := &OkxApproveTransactionResponse{}
	json.Unmarshal(body, out)
	return out, nil
}

func (p *OkxProvider) apiSwapQuote(ctx context.Context, in *SwapQuoteInput) (*OkxSwapQuoteResponse, error) {
	api := fmt.Sprintf("%s/dex/aggregator/quote?chainId=%d&amount=%s&fromTokenAddress=%s&toTokenAddress=%s", OkxProviderAPIBaseURL, in.ChainID, in.Amount, in.FromTokenAddress, in.ToTokenAddress)
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
			return nil, errors.New("Request timeout")
		}
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	out := &OkxSwapQuoteResponse{}
	json.Unmarshal(body, out)
	return out, nil
}

func (p *OkxProvider) apiSwapTrade(ctx context.Context, in *SwapTradeInput) (*OkxSwapTradeResponse, error) {
	var slippage float32
	slippage = in.Slippage * p.config.slippageDecimals
	if slippage < p.config.minSlippage {
		slippage = p.config.minSlippage
	}
	if slippage > p.config.maxSlippage {
		slippage = p.config.maxSlippage
	}
	slippageStr := fmt.Sprintf("%.5f", slippage)
	api := fmt.Sprintf("%s/dex/aggregator/swap?chainId=%d&amount=%s&fromTokenAddress=%s&toTokenAddress=%s&slippage=%s&userWalletAddress=%s", OkxProviderAPIBaseURL, in.ChainID, in.Amount, in.FromTokenAddress, in.ToTokenAddress, slippageStr, in.DestReceiver)
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
	out := &OkxSwapTradeResponse{}
	json.Unmarshal(body, out)
	return out, nil
}

//----------------------------------------------------------------
// api end
//----------------------------------------------------------------
