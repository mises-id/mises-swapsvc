package swap

import (
	"context"

	"github.com/mises-id/mises-swapsvc/app/models"
	"github.com/mises-id/mises-swapsvc/app/models/search"
	"github.com/mises-id/mises-swapsvc/lib/pagination"
)

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
	SwapQuoteInput struct {
		ChainID          uint64
		FromTokenAddress string
		ToTokenAddress   string
		Amount           string
	}

	SwapController struct {
		providers map[string]Provider
	}
)

var (
	swapCtrlEntity *SwapController
)

func NewSwapController() *SwapController {
	if swapCtrlEntity != nil {
		return swapCtrlEntity
	}
	//register provider
	providers := make(map[string]Provider, 0)
	oneProvider := NewOneInchProvider()
	if oneProvider != nil {
		providers[oneProvider.Key()] = oneProvider
	}
	swapCtrlEntity = &SwapController{providers: providers}
	return swapCtrlEntity
}

func GetSwapApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error) {
	ctrl := NewSwapController()
	return ctrl.getApproveAllowance(ctx, in)
}

func ApproveSwapTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error) {
	ctrl := NewSwapController()
	return ctrl.approveTransaction(ctx, in)
}

func SwapTrade(ctx context.Context, in *SwapTradeInput) (*SwapTradeInfo, error) {
	ctrl := NewSwapController()
	return ctrl.swapTrade(ctx, in)
}
func SwapQuote(ctx context.Context, in *SwapQuoteInput) ([]*SwapQuoteInfo, error) {
	ctrl := NewSwapController()
	return ctrl.swapQuote(ctx, in)
}

// PageWebsite
func PageSwapOrder(ctx context.Context, in *SwapOrderInput) ([]*models.SwapOrder, pagination.Pagination, error) {
	params := in.SwapOrderSearch
	return models.PageSwapOrder(ctx, params)
}

// FindSwapOrder
func FindSwapOrder(ctx context.Context, in *SwapOrderInput) (*models.SwapOrder, error) {
	params := in.SwapOrderSearch
	return models.FindSwapOrder(ctx, params)
}

// SwapTokenList
func ListSwapToken(ctx context.Context, in *SwapTokenInput) ([]*models.SwapToken, error) {
	params := in.SwapTokenSearch
	return models.ListSwapToken(ctx, params)
}
