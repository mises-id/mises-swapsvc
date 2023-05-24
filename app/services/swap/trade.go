package swap

import (
	"context"
	"strings"

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

func (c *SwapController) swapTrade(ctx context.Context, in *SwapTradeInput) (*SwapTradeInfo, error) {
	//preload input parameters
	if err := preloadSwapTradeInput(in); err != nil {
		return nil, err
	}
	provider, err := c.findProviderByChainIDAndAddress(ctx, in.ChainID, in.AggregatorAddress)
	if err != nil {
		return nil, err
	}
	return provider.SwapTrade(ctx, in)
}

// ----------------------------------------------------------------
func preloadSwapTradeInput(in *SwapTradeInput) error {
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
	in.AggregatorAddress = strings.ToLower(in.AggregatorAddress)
	return nil
}
