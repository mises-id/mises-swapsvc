package swap

import (
	"context"
	"fmt"
	"time"

	"github.com/mises-id/mises-swapsvc/app/models"
	"github.com/mises-id/mises-swapsvc/app/models/search"
)

type (
	Provider interface {
		GetSwapApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error)
		ApproveSwapTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error)
		SwapQuote(ctx context.Context, in *SwapQuoteInput) *SwapQuoteInfo
		SwapQuotes(ctx context.Context, in *SwapQuoteInput) []*SwapQuoteInfo
		SwapTrade(ctx context.Context, in *SwapTradeInput) (*SwapTradeInfo, error)
		Key() string
	}
	BaseProvider struct {
		swapContractMap map[uint64]*models.SwapContract
		providerInfo    *models.SwapProvider
		ctx             context.Context
		key             string
	}
)

func NewBaseProvider(providerkey string) *BaseProvider {
	resp := &BaseProvider{
		key:             providerkey,
		ctx:             context.TODO(),
		swapContractMap: map[uint64]*models.SwapContract{},
	}
	resp.init()
	return resp
}

func (p *BaseProvider) init() {
	//providerInfo
	p.updateProviderInfo()
	//swapContractMap
	p.updateSwapContractMap()
}

func (p *BaseProvider) updateProviderInfo() error {
	provider, _ := models.FindSwapProivderByKey(p.ctx, p.key)
	p.providerInfo = provider
	return nil
}

func (p *BaseProvider) updateSwapContractMap() error {
	params := &search.SwapContractSearch{ProviderKey: p.key}
	lists, err := models.ListSwapContract(p.ctx, params)
	if err != nil {
		fmt.Printf("[%s] OneInchProvider ListSwapContract error: %s\n", time.Now().Local().String(), err.Error())
		return err
	}
	for _, contract := range lists {
		address := contract.Address
		if address == "" {
			continue
		}
		p.swapContractMap[contract.ChainID] = contract
	}
	return nil
}

func (p *BaseProvider) Key() string {
	return p.key
}

func (p *BaseProvider) findContractByChainID(chainID uint64) *models.SwapContract {
	if p.swapContractMap == nil || len(p.swapContractMap) == 0 {
		p.updateSwapContractMap()
	}
	if p.swapContractMap == nil {
		return nil
	}
	contract, ok := p.swapContractMap[chainID]
	if !ok {
		return nil
	}
	return contract
}

func (p *BaseProvider) aggregator(chainID uint64) *Aggregator {
	if p.providerInfo == nil {
		p.updateProviderInfo()
	}
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
