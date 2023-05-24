package swap

import "context"

type (
	Provider interface {
		GetSwapApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error)
		ApproveSwapTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error)
		SwapQuote(ctx context.Context, in *SwapQuoteInput) *SwapQuoteInfo
		SwapTrade(ctx context.Context, in *SwapTradeInput) (*SwapTradeInfo, error)
		Key() string
	}
)

const (
	ErrUnspportChainID = "Unsupported Chain ID"
)
