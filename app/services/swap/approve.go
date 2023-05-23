package swap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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
	oneInchProvider            = "1inch"
	oKXProvider                = "okx"
	oneIncheProviderAPIBaseURL = "https://api-mises.1inch.io/v5.0"
)

func getApproveAllowance(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error) {
	//check input parameters
	if err := checkApproveAllowanceInput(in); err != nil {
		return nil, err
	}
	//find provider swap contract by chainID & aggregatorAddress
	swapContract, err := findSwapContractByChainIDAndAddress(ctx, in.ChainID, in.AggregatorAddress)
	if err != nil {
		return nil, err
	}
	//get approve allowance by provider

	return getApproveAllowanceByProviderKey(ctx, in, swapContract.ProviderKey)
}

func approveTransaction(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error) {
	//check input parameters
	if err := checkApproveTransactionInput(in); err != nil {
		return nil, err
	}
	//find provider swap contract by chainID & aggregatorAddress
	swapContract, err := findSwapContractByChainIDAndAddress(ctx, in.ChainID, in.AggregatorAddress)
	if err != nil {
		return nil, err
	}
	return approveTransactionByProviderKey(ctx, in, swapContract.ProviderKey)
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
	return nil
}
func approveTransactionByProviderKey(ctx context.Context, in *ApproveSwapTransactionInput, providerKey string) (*ApproveSwapTransactionOutput, error) {
	switch providerKey {
	case oneInchProvider:
		return approveTransactionByOneInchProvider(ctx, in)
	default:
		return nil, codes.ErrInvalidArgument.New("Unsupported aggregatorAddress")
	}
}
func approveTransactionByOneInchProvider(ctx context.Context, in *ApproveSwapTransactionInput) (*ApproveSwapTransactionOutput, error) {
	data, err := apiApproveTransactionByOneInch(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &ApproveSwapTransactionOutput{
		Data:     data.Data,
		GasPrice: data.GasPrice,
		To:       data.To,
		Value:    data.Value,
	}
	return res, nil
}

func apiApproveTransactionByOneInch(ctx context.Context, in *ApproveSwapTransactionInput) (*oneIncheApproveTransactionResponse, error) {
	api := fmt.Sprintf("%s/%d/approve/transaction?tokenAddress=%s&amount=%s", oneIncheProviderAPIBaseURL, in.ChainID, in.TokenAddress, in.Amount)
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	client.Timeout = time.Second * 3
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	out := &oneIncheApproveTransactionResponse{}
	json.Unmarshal(body, out)
	return out, nil
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
	return nil
}

func findSwapContractByChainIDAndAddress(ctx context.Context, chainID uint64, address string) (*models.SwapContract, error) {
	params := &search.SwapContractSearch{
		ChainID: chainID,
		Address: address,
	}
	return models.FindSwapContract(ctx, params)
}

func getApproveAllowanceByProviderKey(ctx context.Context, in *GetSwapApproveAllowanceInput, providerKey string) (*GetSwapApproveAllowanceOutput, error) {
	switch providerKey {
	case oneInchProvider:
		return getApproveAllowanceByOneInchProvider(ctx, in)
	default:
		return nil, codes.ErrInvalidArgument.New("Unsupported aggregatorAddress")
	}
}

func getApproveAllowanceByOneInchProvider(ctx context.Context, in *GetSwapApproveAllowanceInput) (*GetSwapApproveAllowanceOutput, error) {
	data, err := apiGetApproveAllowanceByOneInch(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &GetSwapApproveAllowanceOutput{
		Allowance: data.Allowance,
	}
	return res, nil
}

func apiGetApproveAllowanceByOneInch(ctx context.Context, in *GetSwapApproveAllowanceInput) (*oneIncheGetApproveAllowanceResponse, error) {
	api := fmt.Sprintf("%s/%d/approve/allowance?tokenAddress=%s&walletAddress=%s", oneIncheProviderAPIBaseURL, in.ChainID, in.TokenAddress, in.WalletAddress)
	transport := &http.Transport{Proxy: setProxy()}
	client := &http.Client{Transport: transport}
	client.Timeout = time.Second * 3
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	out := &oneIncheGetApproveAllowanceResponse{}
	json.Unmarshal(body, out)
	return out, nil
}
