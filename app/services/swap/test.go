package swap

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/mises-id/mises-swapsvc/contracts/router"
	"github.com/mises-id/mises-swapsvc/contracts/store"
)

func Test(ctx context.Context) error {
	client, err := ethclient.Dial("https://uk.rpc.blxrbdn.com")
	if err != nil {
		return err
	}
	address := common.HexToAddress("0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6")
	instance, err := store.NewStore(address, client)
	if err != nil {
		return err
	}
	opts := &bind.CallOpts{
		From: common.Address{},
	}
	inTokenAddress := common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
	outTokenAddress := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	raw := &store.StoreRaw{
		Contract: instance,
	}
	fee := big.NewInt(3000)
	amountIn := big.NewInt(200000000)
	sqrtPriceLimitX96 := big.NewInt(0)
	var out []interface{}
	/* var IQuoterV2QuoteExactInputSingleParams *router.ISwapRouterExactInputSingleParams
	IQuoterV2QuoteExactInputSingleParams = &router.ISwapRouterExactInputSingleParams{
		TokenIn:           inTokenAddress,
		TokenOut:          outTokenAddress,
		Fee:               fee,
		AmountIn:          amountIn,
		SqrtPriceLimitX96: sqrtPriceLimitX96,
	} */
	//res, err := instance.QuoteExactInputSingle(&bind.TransactOpts{}, inTokenAddress, outTokenAddress, fee, amountIn, sqrtPriceLimitX96)
	err = raw.Call(opts, &out, "quoteExactInputSingle", inTokenAddress, outTokenAddress, fee, amountIn, sqrtPriceLimitX96)
	if err != nil {
		return err
	}

	fmt.Println("out: ", out) // "1.0"
	return nil
}
