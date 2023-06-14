package swap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type (
	OneInchProviderConfig struct {
		minSlippage, maxSlippage float32
		slippageDecimals         float32
	}
	OneInchProvider struct {
		*BaseProvider
		config    *OneInchProviderConfig
		protocols map[uint64][]*OneInchProtocols
	}
	OneInchProtocols struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		ImgColor  string `json:"img_color"`
		Addresses string `json:"addresses"`
	}
	OneInchProtocolsResponse struct {
		Protocols []*OneInchProtocols `json:"protocols"`
	}
)

func NewOneInchProvider() *OneInchProvider {
	resp := &OneInchProvider{
		BaseProvider: NewBaseProvider(oneInchProviderKey),
		protocols:    map[uint64][]*OneInchProtocols{},
	}
	config := &OneInchProviderConfig{
		minSlippage:      0,
		maxSlippage:      50,
		slippageDecimals: 100,
	}
	resp.config = config
	return resp
}

// GetSwapApproveAllowance
func (p *OneInchProvider) GetSwapApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error) {
	return p.getApproveAllowance(ctx, in)
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

// ApproveSwapTransaction
func (p *OneInchProvider) ApproveSwapTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error) {
	return p.approveTransaction(ctx, in)
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

// SwapQuotes
func (p *OneInchProvider) SwapQuotes(ctx context.Context, in *SwapQuoteInput) []*SwapQuoteInfo {
	if in == nil {
		return nil
	}
	var maxProtocolNum int
	maxProtocolNum = 10
	baseQuoteProtocol := &OneInchProtocols{
		ID: "",
	}
	quoteProtocolList := []*OneInchProtocols{baseQuoteProtocol}
	allProtocol := p.getProtocolListByChainID(in.ChainID)
	totalProtocol := len(allProtocol)
	validProtocol := allProtocol
	if totalProtocol > maxProtocolNum {
		validProtocol = allProtocol[totalProtocol-maxProtocolNum:]
	}
	quoteProtocolList = append(quoteProtocolList, validProtocol...)
	quotes := make([]*SwapQuoteInfo, 0)
	wg := &sync.WaitGroup{}
	checkTask := make(chan int)
	totalTask := len(quoteProtocolList)
	if totalTask == 0 {
		return nil
	}
	wg.Add(1)
	go p.checkTask(totalTask, checkTask, wg)
	for _, quoteProtocol := range quoteProtocolList {
		go func(protocol *OneInchProtocols) {
			wg.Add(1)
			defer func() {
				wg.Done()
				checkTask <- 1
			}()
			params := &OneIncheSwapQuoteParams{
				ChainID:          in.ChainID,
				Amount:           in.Amount,
				FromTokenAddress: in.FromTokenAddress,
				ToTokenAddress:   in.ToTokenAddress,
				Protocols:        protocol.ID,
				SwapFee:          swapFee,
			}
			quote := p.swapQuote(ctx, params)
			if protocol.ID != "" && quote != nil && quote.Aggregator != nil {
				quote.Aggregator.Name = protocol.Title
				quote.Aggregator.Logo = protocol.ImgColor
			}
			quotes = append(quotes, quote)
		}(quoteProtocol)
	}
	wg.Wait()
	return quotes
}

func (c *OneInchProvider) checkTask(taskNum int, ck chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	if taskNum == 0 {
		return
	}
	for {
		<-ck
		taskNum--
		if taskNum == 0 {
			break
		}
	}
}

// SwapQuote
func (p *OneInchProvider) SwapQuote(ctx context.Context, in *SwapQuoteInput) *SwapQuoteInfo {
	if in == nil {
		return nil
	}
	params := &OneIncheSwapQuoteParams{
		ChainID:          in.ChainID,
		Amount:           in.Amount,
		FromTokenAddress: in.FromTokenAddress,
		ToTokenAddress:   in.ToTokenAddress,
		Protocols:        "", //all
		SwapFee:          swapFee,
	}
	return p.swapQuote(ctx, params)
}

// quote
func (p *OneInchProvider) swapQuote(ctx context.Context, in *OneIncheSwapQuoteParams) *SwapQuoteInfo {
	if in == nil {
		return nil
	}
	//check if OneInch is supported. chainId
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
	data, err := p.apiSwapQuoteByOneInchProvider(ctx, in)
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

// SwapTrade
func (p *OneInchProvider) SwapTrade(ctx context.Context, in *SwapTradeInput) (*SwapTradeInfo, error) {
	return p.swapTrade(ctx, in), nil
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
	data, err := p.apiSwapTradeByOneInchProvider(ctx, in)
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

func (p *OneInchProvider) apiSwapTradeByOneInchProvider(ctx context.Context, in *SwapTradeInput) (*oneInchSwapResponse, error) {
	address := strings.ToLower(swapReferrerAddress)
	if !strings.HasPrefix(address, "0x") {
		address = "0x" + address
	}
	var slippage float32
	slippage = in.Slippage * p.config.slippageDecimals
	if slippage < p.config.minSlippage {
		slippage = p.config.minSlippage
	}
	if slippage > p.config.maxSlippage {
		slippage = p.config.maxSlippage
	}
	protocols := ""
	api := fmt.Sprintf("%s/%d/swap?fromTokenAddress=%s&toTokenAddress=%s&amount=%s&fromAddress=%s&destReceiver=%s&slippage=%.3f&referrerAddress=%s&fee=%.3f&protocols=%s", oneIncheProviderAPIBaseURL, in.ChainID, in.FromTokenAddress, in.ToTokenAddress, in.Amount, in.FromAddress, in.DestReceiver, slippage, address, swapFee, protocols)
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

func (p *OneInchProvider) apiSwapQuoteByOneInchProvider(ctx context.Context, in *OneIncheSwapQuoteParams) (*oneInchQuoteResponse, error) {
	protocols := in.Protocols
	api := fmt.Sprintf("%s/%d/quote?fromTokenAddress=%s&toTokenAddress=%s&amount=%s&fee=%.3f&protocols=%s", oneIncheProviderAPIBaseURL, in.ChainID, in.FromTokenAddress, in.ToTokenAddress, in.Amount, in.SwapFee, protocols)
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

func (p *OneInchProvider) getProtocolsByChainID(chainID uint64) (protocols string) {
	chainProtocols := p.getProtocolListByChainID(chainID)
	excludes := map[string]string{}
	switch chainID {
	}
	for _, protocol := range chainProtocols {
		_, isExclude := excludes[protocol.ID]
		if isExclude {
			continue
		}
		if protocols == "" {
			protocols = protocol.ID
		} else {
			protocols += "," + protocol.ID
		}
	}
	return protocols
}

func (p *OneInchProvider) getProtocolListByChainID(chainID uint64) []*OneInchProtocols {
	chainProtocols, ok := p.protocols[chainID]
	if !ok {
		protocolResp, _ := p.apiGetProtocolsByChainID(chainID)
		chainProtocols = make([]*OneInchProtocols, 0)
		if protocolResp != nil {
			chainProtocols = protocolResp.Protocols
			p.protocols[chainID] = chainProtocols
		}
	}
	return chainProtocols
}

func (p *OneInchProvider) apiGetProtocolsByChainID(chainID uint64) (*OneInchProtocolsResponse, error) {
	api := fmt.Sprintf("%s/%d/liquidity-sources", oneIncheProviderAPIBaseURL, chainID)
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
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	out := &OneInchProtocolsResponse{}
	json.Unmarshal(body, out)
	return out, nil
}
