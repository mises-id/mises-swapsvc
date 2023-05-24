package swap

import (
	"context"
	"strings"

	"github.com/mises-id/mises-swapsvc/app/models"
	"github.com/mises-id/mises-swapsvc/app/models/search"
	"github.com/mises-id/mises-swapsvc/lib/codes"
)

type (
	oneIncheGetApproveAllowanceResponse struct {
		Allowance string `json:"allowance"`
	}
	oneIncheApproveTransactionResponse struct {
		Data     string `json:"data"`
		To       string `json:"to"`
		GasPrice string `json:"gasPrice"`
		Value    string `json:"value"`
	}
)

const (
	oneInchProviderKey         = "1inch"
	oKXProvider                = "okx"
	oneIncheProviderAPIBaseURL = "https://api-mises.1inch.io/v5.0"
)

func (c *SwapController) getApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error) {
	//check input parameters
	if err := checkApproveAllowanceInput(in); err != nil {
		return nil, err
	}
	provider, err := c.findProviderByChainIDAndAddress(ctx, in.ChainID, in.AggregatorAddress)
	if err != nil {
		return nil, err
	}
	return provider.GetSwapApproveAllowance(ctx, in)
}

func (c *SwapController) findProviderByChainIDAndAddress(ctx context.Context, chainID uint64, address string) (Provider, error) {
	//find provider swap contract by chainID & aggregatorAddress
	swapContract, err := findSwapContractByChainIDAndAddress(ctx, chainID, address)
	if err != nil {
		return nil, err
	}
	provider, ok := c.providers[swapContract.ProviderKey]
	if !ok || provider == nil {
		return nil, codes.ErrInvalidArgument.New("Unsupported aggregatorAddress")
	}
	return provider, nil
}
func (c *SwapController) approveTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error) {
	// check input parameters
	if err := checkApproveTransactionInput(in); err != nil {
		return nil, err
	}
	provider, err := c.findProviderByChainIDAndAddress(ctx, in.ChainID, in.AggregatorAddress)
	if err != nil {
		return nil, err
	}
	return provider.ApproveSwapTransaction(ctx, in)
}

// ----------------------------------------------------------------
// approve transaction
func checkApproveTransactionInput(in *ApproveSwapTransactionInput) error {
	if in.ChainID == 0 {
		return codes.ErrInvalidArgument.New("Invaild chainID")
	}
	if in.AggregatorAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild aggregatorAddress")
	}
	if in.TokenAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild tokenAddress")
	}
	in.AggregatorAddress = strings.ToLower(in.AggregatorAddress)
	return nil
}

// ----------------------------------------------------------------
// approve allowance
func checkApproveAllowanceInput(in *GetSwapApproveAllowanceInput) error {
	if in.ChainID == 0 {
		return codes.ErrInvalidArgument.New("Invaild chainID")
	}
	if in.AggregatorAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild aggregatorAddress")
	}
	if in.WalletAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild walletAddress")
	}
	if in.TokenAddress == "" {
		return codes.ErrInvalidArgument.New("Invaild tokenAddress")
	}
	in.AggregatorAddress = strings.ToLower(in.AggregatorAddress)
	return nil
}

func findSwapContractByChainIDAndAddress(ctx context.Context, chainID uint64, address string) (*models.SwapContract, error) {
	params := &search.SwapContractSearch{
		ChainID: chainID,
		Address: address,
	}
	return models.FindSwapContract(ctx, params)
}
