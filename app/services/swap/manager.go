package swap

import (
	"context"

	"github.com/mises-id/mises-swapsvc/app/services/swap/swap_sync"
	"github.com/mises-id/mises-swapsvc/config/env"
)

type ()

const (
	SyncSwapOrderOpen = "open"
)

var (
	syncSwapOrderFlag = env.Envs.SyncSwapOrderMode
)

// SyncSwapOrder
func SyncSwapOrder(ctx context.Context) error {
	if syncSwapOrderFlag != SyncSwapOrderOpen {
		return nil
	}
	swapSync := swap_sync.NewSwapSync()
	err := swapSync.Run()
	if err != nil {
		return err
	}
	return nil
}

func InitSwap(ctx context.Context) error {
	/* UpdateSwapChain(ctx)
	UpdateSwapContract(ctx)
	UpdateSwapToken(ctx)
	UpdateSwapProvider(ctx)
	models.EnsureIndex() */
	return nil
}
