package swap

import "github.com/mises-id/mises-swapsvc/app/models/search"

type (
	SwapOrderInput struct {
		*search.SwapOrderSearch
	}
	SwapTokenInput struct {
		*search.SwapTokenSearch
	}
	GetSwapApproveAllowanceInput struct {
		ChainID           uint64
		AggregatorAddress string
		TokenAddress      string
		WalletAddress     string
	}
	GetSwapApproveAllowanceOutput struct {
		Allowance string
	}
	ApproveSwapTransactionInput struct {
		ChainID           uint64
		AggregatorAddress string
		TokenAddress      string
		Amount            string
	}
	ApproveSwapTransactionOutput struct {
		Data, To, GasPrice, Value string
	}
	SwapTradeInput struct {
		ChainID           uint64
		FromAddress       string
		FromTokenAddress  string
		ToTokenAddress    string
		Amount            string
		Slippage          float32
		DestReceiver      string
		AggregatorAddress string
	}
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
		ComparePercent   float32
	}
	SwapQuoteInput struct {
		ChainID          uint64
		FromTokenAddress string
		ToTokenAddress   string
		Amount           string
	}
	SwapQuoteOutput struct {
		BestQuote *SwapQuoteInfo
		Error     string `json:"error"`
		AllQuote  []*SwapQuoteInfo
	}
	SwapController struct {
		providers map[string]Provider
	}
	OneIncheSwapQuoteParams struct {
		ChainID          uint64
		FromTokenAddress string
		ToTokenAddress   string
		Amount           string
		Protocols        string
		SwapFee          float32
	}
)
