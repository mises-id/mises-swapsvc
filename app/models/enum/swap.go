package enum

type (
	SwapOrderStatus int
	ChainID uint64
)

const (
	SwapOrderPending SwapOrderStatus = 3
	SwapOrderFail SwapOrderStatus = 2
	SwapOrderSuccess SwapOrderStatus = 1
	//ChainID
	ChainETH ChainID = 1
	ChainBSC ChainID = 56
)