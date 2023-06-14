package swap

import (
	"context"
	"math/big"
	"sort"
	"strconv"
	"sync"

	"github.com/mises-id/mises-swapsvc/lib/codes"
)

type (
	oneInchQuoteResponse struct {
		StatusCode     string     `json:"statusCode"`
		Error          string     `json:"error"`
		Description    string     `json:"description"`
		ToTokenAmount  string     `json:"toTokenAmount"`
		EstimateGasFee int64      `json:"estimatedGas"`
		Tx             *oneInchTx `json:"tx"`
	}
)

func (c *SwapController) swapQuote(ctx context.Context, in *SwapQuoteInput) (*SwapQuoteOutput, error) {
	// check input parameters
	if err := checkSwapQuoteInput(in); err != nil {
		return nil, err
	}
	// get quotes
	quotes, err := c.getSwapQuotes(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &SwapQuoteOutput{}
	if len(quotes) == 0 {
		resp.Error = ErrorNotFoundQuote
	}
	validQuotes, bestQuote := c.getValidQuotesAndBestQuote(ctx, quotes)
	resp.AllQuote = validQuotes
	resp.BestQuote = bestQuote
	return resp, nil
}

func (c *SwapController) getSwapQuotes(ctx context.Context, in *SwapQuoteInput) ([]*SwapQuoteInfo, error) {
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
			quote := provider.SwapQuotes(ctx, in)
			quotes = append(quotes, quote...)
		}(provider)
	}
	wg.Wait()
	return quotes, nil
}

func (c *SwapController) getValidQuotesAndBestQuote(ctx context.Context, quotes []*SwapQuoteInfo) (validQuotes []*SwapQuoteInfo, bestQuote *SwapQuoteInfo) {
	if quotes == nil || len(quotes) == 0 {
		return nil, nil
	}
	var worstQuote *SwapQuoteInfo
	validQuotes = make([]*SwapQuoteInfo, 0)
	//get best quote
	for _, quote := range quotes {
		//check
		if quote.Error != "" || quote.ToTokenAmount == "" {
			continue
		}
		validQuotes = append(validQuotes, quote)
		toTokenAmount, _ := new(big.Int).SetString(quote.ToTokenAmount, 10)
		if bestQuote == nil {
			bestQuote = quote
			worstQuote = quote
			continue
		}
		tempToTokenAmount, _ := new(big.Int).SetString(bestQuote.ToTokenAmount, 10)
		if toTokenAmount.Cmp(tempToTokenAmount) == 1 {
			bestQuote = quote
		}
		if toTokenAmount.Cmp(tempToTokenAmount) == -1 {
			worstQuote = quote
		}
	}
	//best compare percent
	if bestQuote != nil {
		bestQuote.ComparePercent = compareQuotePercent(bestQuote, worstQuote)
	}
	//valid quotes
	for _, quote := range quotes {
		if quote.Error != "" || quote.ToTokenAmount == "" {
			quote.ComparePercent = -100
			continue
		}
		if quote != bestQuote {
			quote.ComparePercent = compareQuotePercent(quote, bestQuote)
		}
	}
	sort.Slice(quotes, func(i, j int) bool {
		return quotes[i].ComparePercent > quotes[j].ComparePercent
	})
	return quotes, bestQuote
}

func compareQuotePercent(source *SwapQuoteInfo, dest *SwapQuoteInfo) (comparePercent float32) {
	if source != nil && dest != nil && (dest.ToTokenAmount != "" && dest.ToTokenAmount != "0") && (source.ToTokenAmount != "" && source.ToTokenAmount != "0") {
		bestToTokenAmount, _ := new(big.Float).SetString(source.ToTokenAmount)
		worstToTokenAmount, _ := new(big.Float).SetString(dest.ToTokenAmount)
		if source.ToTokenAmount != dest.ToTokenAmount {
			var overAmount, comparePercentBig *big.Float
			overAmount = big.NewFloat(10)
			overAmount = overAmount.Sub(bestToTokenAmount, worstToTokenAmount)
			comparePercentBig = overAmount.Quo(overAmount, worstToTokenAmount)
			comparePercentStr := comparePercentBig.String()
			comparePercent64, _ := strconv.ParseFloat(comparePercentStr, 10)
			comparePercent = float32(comparePercent64) * 100
		}
	}
	return comparePercent
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
