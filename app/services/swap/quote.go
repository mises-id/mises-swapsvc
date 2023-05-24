package swap

import (
	"context"
	"sync"

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

func (c *SwapController) swapQuote(ctx context.Context, in *SwapQuoteInput) ([]*SwapQuoteInfo, error) {
	// check input parameters
	if err := checkSwapQuoteInput(in); err != nil {
		return nil, err
	}
	wg := &sync.WaitGroup{}
	checkTask := make(chan int)
	totalTask := len(c.providers)
	if totalTask == 0 {
		return nil, nil
	}
	wg.Add(1)
	go c.checkTask(totalTask, checkTask, wg)
	quotes := make([]*SwapQuoteInfo, 0)
	for _, provider := range c.providers {
		go func(provider Provider) {
			wg.Add(1)
			defer func() {
				wg.Done()
				checkTask <- 1
			}()
			//run
			quote := provider.SwapQuote(ctx, in)
			quotes = append(quotes, quote)
		}(provider)
	}
	wg.Wait()
	return quotes, nil
}

func (c *SwapController) checkTask(taskNum int, ck chan int, wg *sync.WaitGroup) {
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

// ----------------------------------------------------------------
func checkSwapQuoteInput(in *SwapQuoteInput) error {
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
