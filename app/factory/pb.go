package factory

import (
	"strconv"

	"github.com/mises-id/mises-swapsvc/app/models"
	pb "github.com/mises-id/mises-swapsvc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func docID(id primitive.ObjectID) string {
	if id.IsZero() {
		return ""
	}
	return id.Hex()
}

func NewSwapTokenlice(data []*models.SwapToken) []*pb.Token {
	result := make([]*pb.Token, len(data))
	for i, v := range data {
		result[i] = NewSwapToken(v)
	}
	return result
}

func NewSwapToken(data *models.SwapToken) *pb.Token {
	if data == nil {
		return nil
	}
	resp := pb.Token{
		Address:  data.Address,
		Name:     data.Name,
		LogoUri:  data.LogoURI,
		Decimals: int32(data.Decimals),
		Symbol:   data.Symbol,
		ChainID:  data.ChainID,
	}
	return &resp
}
func NewSwapOrderSlice(data []*models.SwapOrder) []*pb.SwapOrder {
	result := make([]*pb.SwapOrder, len(data))
	for i, v := range data {
		result[i] = NewSwapOrder(v)
	}
	return result
}

func NewSwapOrder(data *models.SwapOrder) *pb.SwapOrder {
	if data == nil {
		return nil
	}
	resp := pb.SwapOrder{
		Id:              docID(data.ID),
		ChainID:         data.ChainID,
		FromAddress:     data.FromAddress,
		DestReceiver:    data.DestReceiver,
		ReceiptStatus:   int32(data.ReceiptStatus),
		ContractAddress: data.ContractAddress,
		TransactionFee:  data.TransactionFee,
		BlockAt:         data.BlockAt.Unix(),
	}
	if data.Transaction != nil {
		resp.Tx = NewTransaction(data.Transaction)
	}
	if data.Provider != nil {
		resp.Provider = NewSwapProvider(data.Provider)
	}
	if data.FromToken != nil {
		resp.FromToken = NewFromToken(data.FromToken)
	}
	if data.ToToken != nil {
		resp.ToToken = NewToToken(data.ToToken)
	}
	resp.NativeToken = NewToken(data.NativeToken)
	return &resp
}

func NewSwapProvider(data *models.SwapProvider) *pb.SwapProvider {
	if data == nil {
		return nil
	}
	resp := pb.SwapProvider{
		Key:  data.Key,
		Name: data.Name,
		Logo: data.Logo,
	}
	return &resp
}
func NewTransaction(data *models.Transaction) *pb.Transaction {
	if data == nil {
		return nil
	}
	blockNumber, _ := strconv.ParseInt(data.BlockNumber, 10, 64)
	resp := pb.Transaction{
		Hash:        data.Hash,
		Gas:         data.Gas,
		BlockNumber: blockNumber,
		GasUsed:     data.GasUsed,
		GasPrice:    data.GasPrice,
		Nonce:       data.Nonce,
	}
	return &resp
}

func NewToken(data *models.Token) *pb.Token {
	if data == nil {
		return nil
	}
	resp := pb.Token{
		Address:  data.Address,
		Name:     data.Name,
		LogoUri:  data.LogoURI,
		Decimals: int32(data.Decimals),
		Symbol:   data.Symbol,
	}
	return &resp
}
func NewFromToken(data *models.FromToken) *pb.Token {
	if data == nil {
		return nil
	}
	resp := pb.Token{
		Address: data.Address,
		Value:   data.Value,
	}
	if data.Token != nil {
		resp.Name = data.Token.Name
		resp.LogoUri = data.Token.LogoURI
		resp.Decimals = int32(data.Token.Decimals)
		resp.Symbol = data.Token.Symbol
	}
	return &resp
}
func NewToToken(data *models.ToToken) *pb.Token {
	if data == nil {
		return nil
	}
	resp := pb.Token{
		Address: data.Address,
		Value:   data.Value,
	}
	if data.Token != nil {
		resp.Name = data.Token.Name
		resp.LogoUri = data.Token.LogoURI
		resp.Decimals = int32(data.Token.Decimals)
		resp.Symbol = data.Token.Symbol
	}
	return &resp
}
