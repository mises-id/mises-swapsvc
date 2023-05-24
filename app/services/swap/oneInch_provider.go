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

	"github.com/mises-id/mises-swapsvc/app/models"
	"github.com/mises-id/mises-swapsvc/app/models/search"
)

type (
	OneInchProvider struct {
		key             string
		swapContractMap map[uint64]*models.SwapContract
		providerInfo    *models.SwapProvider
		ctx             context.Context
	}
)

func NewOneInchProvider() *OneInchProvider {
	resp := &OneInchProvider{
		key:             oneInchProviderKey,
		ctx:             context.TODO(),
		swapContractMap: map[uint64]*models.SwapContract{},
	}
	resp.init()
	return resp
}

func (p *OneInchProvider) init() {
	//providerInfo
	provider, _ := models.FindSwapProivderByKey(p.ctx, p.key)
	p.providerInfo = provider
	//swapContractMap
	params := &search.SwapContractSearch{ProviderKey: p.key}
	lists, err := models.ListSwapContract(p.ctx, params)
	if err != nil {
		fmt.Printf("[%s] OneInchProvider ListSwapContract error: %s\n", time.Now().Local().String(), err.Error())
		return
	}
	for _, contract := range lists {
		address := contract.Address
		if address == "" {
			continue
		}
		p.swapContractMap[contract.ChainID] = contract
	}
}

// GetSwapApproveAllowance
func (p *OneInchProvider) GetSwapApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error) {

	return p.getApproveAllowance(ctx, in)
}

// ApproveSwapTransaction
func (p *OneInchProvider) ApproveSwapTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error) {
	return p.approveTransaction(ctx, in)
}

// SwapQuote
func (p *OneInchProvider) SwapQuote(ctx context.Context, in *SwapQuoteInput) *SwapQuoteInfo {
	return p.swapQuote(ctx, in)
}

// SwapTrade
func (p *OneInchProvider) SwapTrade(ctx context.Context, in *SwapTradeInput) (*SwapTradeInfo, error) {
	return p.swapTrade(ctx, in), nil
}

// Key
func (p *OneInchProvider) Key() string {
	return p.key
}

func (p *OneInchProvider) approveTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error) {
	data, err := apiApproveTransactionByOneInch(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &ApproveSwapTransactionOutput{
		Data:     data.Data,
		GasPrice: data.GasPrice,
		To:       data.To,
		Value:    data.Value,
	}
	return res, nil
}

func apiApproveTransactionByOneInch(ctx context.Context, in *ApproveSwapTransactionInput) (*oneIncheApproveTransactionResponse, error) {
	api := fmt.Sprintf("%s/%d/approve/transaction?tokenAddress=%s&amount=%s", oneIncheProviderAPIBaseURL, in.ChainID, in.TokenAddress, in.Amount)
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
	out := &oneIncheApproveTransactionResponse{}
	json.Unmarshal(body, out)
	return out, nil
}

func (p *OneInchProvider) getApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error) {
	data, err := apiGetApproveAllowanceByOneInch(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &GetSwapApproveAllowanceOutput{
		Allowance: data.Allowance,
	}
	return res, nil
}

// swap
func (p *OneInchProvider) swapTrade(ctx context.Context, in *SwapTradeInput) *SwapTradeInfo {
	if in == nil {
		return nil
	}
	//check if OneInch is supported
	st := time.Now()
	resp := &SwapTradeInfo{
		FromTokenAddress: in.FromTokenAddress,
		ToTokenAddress:   in.ToTokenAddress,
		Fee:              swapFee,
		FromTokenAmount:  in.Amount,
	}
	resp.Aggregator = p.aggregator(in.ChainID)
	contract := p.findContractByChainID(in.ChainID)
	if contract == nil {
		resp.Error = ErrUnspportChainID
		return resp
	}
	data, err := apiSwapTradeByOneInchProvider(ctx, in)
	resp.FetchTime = int64(time.Since(st).Milliseconds())
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	resp.Error = data.Description
	resp.ToTokenAmount = data.ToTokenAmount
	if data.Tx != nil {
		resp.Trade = &Trade{
			From:     data.Tx.From,
			To:       data.Tx.To,
			Data:     data.Tx.Data,
			GasPrice: data.Tx.GasPrice,
			GasLimit: fmt.Sprintf("%d", data.Tx.Gas),
			Value:    data.Tx.Value,
		}
	}
	return resp
}

// quote
func (p *OneInchProvider) swapQuote(ctx context.Context, in *SwapQuoteInput) *SwapQuoteInfo {
	//check if OneInch is supported
	st := time.Now()
	resp := &SwapQuoteInfo{
		FromTokenAddress: in.FromTokenAddress,
		ToTokenAddress:   in.ToTokenAddress,
		Fee:              swapFee,
		FromTokenAmount:  in.Amount,
	}
	resp.Aggregator = p.aggregator(in.ChainID)
	contract := p.findContractByChainID(in.ChainID)
	if contract == nil {
		resp.Error = ErrUnspportChainID
		return resp
	}
	data, err := apiSwapQuoteByOneInchProvider(ctx, in)
	resp.FetchTime = int64(time.Since(st).Milliseconds())
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	resp.EstimateGasFee = fmt.Sprintf("%d", data.EstimateGasFee)
	resp.Error = data.Description
	resp.ToTokenAmount = data.ToTokenAmount
	return resp
}

func apiGetApproveAllowanceByOneInch(ctx context.Context, in *GetSwapApproveAllowanceInput) (*oneIncheGetApproveAllowanceResponse, error) {
	api := fmt.Sprintf("%s/%d/approve/allowance?tokenAddress=%s&walletAddress=%s", oneIncheProviderAPIBaseURL, in.ChainID, in.TokenAddress, in.WalletAddress)
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
	out := &oneIncheGetApproveAllowanceResponse{}
	json.Unmarshal(body, out)
	return out, nil
}

func (p *OneInchProvider) aggregator(chainID uint64) *Aggregator {
	if p.providerInfo == nil {
		return nil
	}
	resp := &Aggregator{
		Logo:            p.providerInfo.Logo,
		Type:            p.providerInfo.Type,
		Name:            p.providerInfo.Name,
		ContractAddress: "",
	}
	contract := p.findContractByChainID(chainID)
	if contract != nil {
		resp.ContractAddress = contract.Address
	}
	return resp
}

func (p *OneInchProvider) findContractByChainID(chainID uint64) *models.SwapContract {
	contract, ok := p.swapContractMap[chainID]
	if !ok {
		return nil
	}
	return contract
}

func apiSwapTradeByOneInchProvider(ctx context.Context, in *SwapTradeInput) (*oneInchSwapResponse, error) {
	address := strings.ToLower(swapReferrerAddress)
	if !strings.HasPrefix(address, "0x") {
		address = "0x" + address
	}
	api := fmt.Sprintf("%s/%d/swap?fromTokenAddress=%s&toTokenAddress=%s&amount=%s&fromAddress=%s&destReceiver=%s&slippage=%.3f&referrerAddress=%s&fee=%.3f", oneIncheProviderAPIBaseURL, in.ChainID, in.FromTokenAddress, in.ToTokenAddress, in.Amount, in.FromAddress, in.DestReceiver, in.Slippage, address, swapFee)
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
	out := &oneInchSwapResponse{}
	json.Unmarshal(body, out)
	return out, nil
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