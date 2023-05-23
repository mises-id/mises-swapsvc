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

	"github.com/mises-id/mises-swapsvc/config/env"
	"github.com/mises-id/mises-swapsvc/lib/codes"
)

type (
	SwapTradeInfo struct {
		Aggregator       *Aggregator `json:"aggregator"`       //聚合器信息
		FromTokenAddress string      `json:"fromTokenAddress"` //swap发起询价的token地址
		ToTokenAddress   string      `json:"toTokenAddress"`   //swap目标token
		FromTokenAmount  string      `json:"fromTokenAmount"`  //swap发起询价token的数目
		ToTokenAmount    string      `json:"toTokenAmount"`    //swap目标token的数目
		Trade            *Trade      `json:"trade"`            //swap交易数据，当Error不等于空时为null
		Fee              float32     `json:"fee"`              //收取的佣金百分比，FromTokenAmount令牌数量的这个百分比将被发送到 referrerAddress，其余部分将用作交换的输入
		Error            string      `json:"error"`
		FetchTime        int64       `json:"fetchTime"`
	}
	Trade struct {
		Data     string `json:"data"`      //Call data
		From     string `json:"from"`      //用户钱包地址
		To       string `json:"to"`        //router 合约地址如：1inchV5
		GasPrice string `json:"gas_price"` //以 wei 为单位的 gas price
		GasLimit string `json:"gas_limit"` //gas limit 的估计量
		Value    string `json:"value"`     //与合约交互的主链币数量（wei)
	}
	Aggregator struct {
		Logo            string `json:"logo"`
		Type            string `json:"type"`             //聚合类型 RFQ，AGG
		Name            string `json:"name"`             //聚合器名字，如1inch v5: Aggregation Router 请求交易授权需要传此值
		ContractAddress string `json:"contract_address"` //聚合器Router合约地址
	}
	Token struct {
		Address  string `json:"address"`  //token 地址
		Decimals int32  `json:"decimals"` //小数位
		LogoUri  string `json:"logo_uri"` //logo
		Name     string `json:"name"`     //名称
		Symbol   string `json:"symbol"`   //token 单位符号
	}
	oneInchTx struct {
		Data     string `json:"data"`
		From     string `json:"from"`
		To       string `json:"to"`
		GasPrice string `json:"gasPrice"`
		Gas      int64  `json:"gas"`
		Value    string `json:"value"`
	}
	oneInchSwapResponse struct {
		StatusCode    string     `json:"statusCode"`
		Error         string     `json:"error"`
		Description   string     `json:"description"`
		ToTokenAmount string     `json:"toTokenAmount"`
		Tx            *oneInchTx `json:"tx"`
	}
)

var (
	swapReferrerAddress = env.Envs.SwapReferrerAddress
	swapFee             float32
	maxSwapFee          float32 = 3
)

func init() {
	swapFee = env.Envs.SwapFee
	if swapFee > maxSwapFee {
		swapFee = maxSwapFee
	}
}

func swapTrade(ctx context.Context, in *SwapTradesInput) (*SwapTradeInfo, error) {
	//preload input parameters
	if err := preloadSwapTradeInput(in); err != nil {
		return nil, err
	}
	//find provider swap contract by chainID & aggregatorAddress
	swapContract, err := findSwapContractByChainIDAndAddress(ctx, in.ChainID, in.AggregatorAddress)
	if err != nil {
		return nil, err
	}
	return getSwapTradeByProvider(ctx, in, swapContract.ProviderKey)
}

func getSwapTradeByProvider(ctx context.Context, in *SwapTradesInput, providerKey string) (*SwapTradeInfo, error) {
	switch providerKey {
	case oneInchProvider:
		return swapTradeByOneInchProvider(ctx, in), nil
	default:
		return nil, codes.ErrInvalidArgument.New("Unsupported aggregatorAddress")
	}
}

func getAggregatorByProviderKey(ctx context.Context, providerKey string) *Aggregator {
	switch providerKey {
	case oneInchProvider:
		return &Aggregator{
			Logo:            "https://cdn.mises.site/s3://mises-storage/upload/website/logo/app_1inch_io_logo703043.io?sign=5XbOzRfkJVskb2A84x46gOnz9h2O1juQC8FRUL07fzg&version=2.0",
			Type:            "AGG",
			Name:            "1inch",
			ContractAddress: "0x1111111254eeb25477b68fb85ed929f73a960582",
		}
	default:
		return nil
	}
}

func swapTrades(ctx context.Context, in *SwapTradesInput) ([]*SwapTradeInfo, error) {
	//preload input parameters
	if err := preloadSwapTradesInput(in); err != nil {
		return nil, err
	}
	//find providers
	//1inch
	trades := make([]*SwapTradeInfo, 0)

	tradeOneInch := swapTradeByOneInchProvider(ctx, in)

	trades = append(trades, tradeOneInch)
	return trades, nil
}

// ----------------------------------------------------------------
func preloadSwapTradeInput(in *SwapTradesInput) error {
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
	if in.FromAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild fromAddress")
	}
	if in.DestReceiver == "" {
		in.DestReceiver = in.FromAddress
	}
	if in.AggregatorAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild aggregatorAddress")
	}
	return nil
}

// ----------------------------------------------------------------
func preloadSwapTradesInput(in *SwapTradesInput) error {
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
	if in.FromAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild fromAddress")
	}
	if in.DestReceiver == "" {
		in.DestReceiver = in.FromAddress
	}
	return nil
}

func swapTradeByOneInchProvider(ctx context.Context, in *SwapTradesInput) *SwapTradeInfo {
	//check if OneInch is supported
	st := time.Now()
	resp := &SwapTradeInfo{
		FromTokenAddress: in.FromTokenAddress,
		ToTokenAddress:   in.ToTokenAddress,
		Fee:              swapFee,
		FromTokenAmount:  in.Amount,
	}
	resp.Aggregator = getAggregatorByProviderKey(ctx, oneInchProvider)
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

func apiSwapTradeByOneInchProvider(ctx context.Context, in *SwapTradesInput) (*oneInchSwapResponse, error) {
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
