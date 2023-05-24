package handlers

import (
	"context"
	"github.com/mises-id/mises-swapsvc/app/factory"
	"github.com/mises-id/mises-swapsvc/app/models/search"
	"github.com/mises-id/mises-swapsvc/app/services/swap"
	"github.com/mises-id/mises-swapsvc/lib/codes"
	"github.com/mises-id/mises-swapsvc/lib/pagination"
	pb "github.com/mises-id/mises-swapsvc/proto"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.SwapsvcServer {
	return swapsvcService{}
}

type swapsvcService struct{}

func (s swapsvcService) SyncSwapOrder(ctx context.Context, in *pb.SyncSwapOrderRequest) (*pb.SyncSwapOrderResponse, error) {
	var resp pb.SyncSwapOrderResponse
	err := swap.SyncSwapOrder(ctx)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s swapsvcService) FindSwapOrder(ctx context.Context, in *pb.FindSwapOrderRequest) (*pb.FindSwapOrderResponse, error) {
	var resp pb.FindSwapOrderResponse
	params := &search.SwapOrderSearch{
		ChainID:     in.ChainID,
		FromAddress: in.FromAddress,
		Hash:        in.TxHash,
	}
	data, err := swap.FindSwapOrder(ctx, &swap.SwapOrderInput{SwapOrderSearch: params})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Data = factory.NewSwapOrder(data)
	return &resp, nil
}

func (s swapsvcService) SwapOrderPage(ctx context.Context, in *pb.SwapOrderPageRequest) (*pb.SwapOrderPageResponse, error) {
	var resp pb.SwapOrderPageResponse
	params := &search.SwapOrderSearch{
		ChainID:     in.ChainID,
		FromAddress: in.FromAddress,
	}
	if in.Paginator != nil {
		params.PageNum = int64(in.Paginator.PageNum)
		params.PageSize = int64(in.Paginator.PageSize)
	}
	data, page, err := swap.PageSwapOrder(ctx, &swap.SwapOrderInput{SwapOrderSearch: params})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Data = factory.NewSwapOrderSlice(data)
	tradpage := page.BuildJSONResult().(*pagination.TraditionalPagination)
	resp.Paginator = &pb.Page{
		PageNum:      uint64(tradpage.PageNum),
		PageSize:     uint64(tradpage.PageSize),
		TotalPage:    uint64(tradpage.TotalPages),
		TotalRecords: uint64(tradpage.TotalRecords),
	}
	return &resp, nil
}

func (s swapsvcService) SwapTrade(ctx context.Context, in *pb.SwapTradeRequest) (*pb.SwapTradeResponse, error) {
	var resp pb.SwapTradeResponse
	params := &swap.SwapTradeInput{
		ChainID:           in.ChainID,
		FromTokenAddress:  in.FromTokenAddress,
		Amount:            in.Amount,
		ToTokenAddress:    in.ToTokenAddress,
		DestReceiver:      in.DestReceiver,
		FromAddress:       in.FromAddress,
		Slippage:          in.Slippage,
		AggregatorAddress: in.AggregatorAddress,
	}
	data, err := swap.SwapTrade(ctx, params)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Data = buildSwapTradeInfo(data)
	return &resp, nil
}

func buildSwapTradeInfo(data *swap.SwapTradeInfo) *pb.SwapTradeInfo {
	if data == nil {
		return nil
	}
	resp := &pb.SwapTradeInfo{
		FromTokenAddress: data.FromTokenAddress,
		ToTokenAddress:   data.ToTokenAddress,
		FromTokenAmount:  data.FromTokenAmount,
		ToTokenAmount:    data.ToTokenAmount,
		Error:            data.Error,
		Fee:              data.Fee,
		FetchTime:        data.FetchTime,
	}
	if data.Aggregator != nil {
		resp.Aggregator = &pb.Aggregator{
			Type:            data.Aggregator.Type,
			Logo:            data.Aggregator.Logo,
			Name:            data.Aggregator.Name,
			ContractAddress: data.Aggregator.ContractAddress,
		}
	}
	if data.Trade != nil {
		resp.Trade = &pb.Trade{
			Data:     data.Trade.Data,
			From:     data.Trade.From,
			To:       data.Trade.To,
			GasPrice: data.Trade.GasPrice,
			GasLimit: data.Trade.GasLimit,
			Value:    data.Trade.Value,
		}
	}
	return resp
}

func (s swapsvcService) ListSwapToken(ctx context.Context, in *pb.ListSwapTokenRequest) (*pb.ListSwapTokenResponse, error) {
	var resp pb.ListSwapTokenResponse
	//check input parameters
	if in.ChainID <= 0 {
		return nil, codes.ErrInvalidArgument.New("chainID required")
	}
	params := &search.SwapTokenSearch{
		ChainID: in.ChainID,
		Address: in.TokenAddress,
	}
	data, err := swap.ListSwapToken(ctx, &swap.SwapTokenInput{SwapTokenSearch: params})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Data = factory.NewSwapTokenlice(data)
	return &resp, nil
}

func (s swapsvcService) GetSwapApproveAllowance(ctx context.Context, in *pb.GetSwapApproveAllowanceRequest) (*pb.GetSwapApproveAllowanceResponse, error) {
	var resp pb.GetSwapApproveAllowanceResponse
	params := &swap.GetSwapApproveAllowanceInput{
		ChainID:           in.ChainID,
		TokenAddress:      in.TokenAddress,
		WalletAddress:     in.WalletAddress,
		AggregatorAddress: in.AggregatorAddress,
	}
	data, err := swap.GetSwapApproveAllowance(ctx, params)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Allowance = data.Allowance
	return &resp, nil
}

func (s swapsvcService) ApproveSwapTransaction(ctx context.Context, in *pb.ApproveSwapTransactionRequest) (*pb.ApproveSwapTransactionResponse, error) {
	var resp pb.ApproveSwapTransactionResponse
	params := &swap.ApproveSwapTransactionInput{
		ChainID:           in.ChainID,
		TokenAddress:      in.TokenAddress,
		Amount:            in.Amount,
		AggregatorAddress: in.AggregatorAddress,
	}
	data, err := swap.ApproveSwapTransaction(ctx, params)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Data = data.Data
	resp.To = data.To
	resp.GasPrice = data.GasPrice
	resp.Value = data.Value
	return &resp, nil
}

func (s swapsvcService) SwapQuote(ctx context.Context, in *pb.SwapQuoteRequest) (*pb.SwapQuoteResponse, error) {
	var resp pb.SwapQuoteResponse
	params := &swap.SwapQuoteInput{
		ChainID:          in.ChainID,
		FromTokenAddress: in.FromTokenAddress,
		Amount:           in.Amount,
		ToTokenAddress:   in.ToTokenAddress,
	}
	data, err := swap.SwapQuote(ctx, params)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Data = buildSwapQuoteSlice(data)
	return &resp, nil
}
func buildSwapQuoteSlice(data []*swap.SwapQuoteInfo) []*pb.SwapQuoteInfo {
	result := make([]*pb.SwapQuoteInfo, len(data))
	for i, v := range data {
		result[i] = buildSwapQuoteInfo(v)
	}
	return result
}

func buildSwapQuoteInfo(data *swap.SwapQuoteInfo) *pb.SwapQuoteInfo {
	if data == nil {
		return nil
	}
	resp := &pb.SwapQuoteInfo{
		FromTokenAddress: data.FromTokenAddress,
		ToTokenAddress:   data.ToTokenAddress,
		FromTokenAmount:  data.FromTokenAmount,
		ToTokenAmount:    data.ToTokenAmount,
		Error:            data.Error,
		Fee:              data.Fee,
		FetchTime:        data.FetchTime,
		EstimateGasFee:   data.EstimateGasFee,
	}
	if data.Aggregator != nil {
		resp.Aggregator = &pb.Aggregator{
			Logo:            data.Aggregator.Logo,
			Type:            data.Aggregator.Type,
			Name:            data.Aggregator.Name,
			ContractAddress: data.Aggregator.ContractAddress,
		}
	}
	return resp
}

func (s swapsvcService) Test(ctx context.Context, in *pb.TestRequest) (*pb.TestResponse, error) {
	var resp pb.TestResponse
	err := swap.InitSwap(ctx)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}
