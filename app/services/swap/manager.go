package swap

import (
	"context"

	"github.com/mises-id/mises-swapsvc/app/services/swap/swap_sync"
)

type ()

// SyncSwapOrder
func SyncSwapOrder(ctx context.Context) error {
	swapSync := swap_sync.NewSwapSync()
	err := swapSync.Run()
	if err != nil {
		return err
	}
	return nil
}
