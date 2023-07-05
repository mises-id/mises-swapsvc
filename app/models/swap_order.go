package models

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/mises-id/mises-swapsvc/app/models/enum"
	"github.com/mises-id/mises-swapsvc/lib/db"
	"github.com/mises-id/mises-swapsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	NativeTokenAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

type (
	Transaction struct {
		Timestamp        string `json:"timestamp"`
		BlockHash        string `json:"blockHash" bson:"blockHash"`
		BlockNumber      string `json:"blockNumber" bson:"blockNumber"`
		From             string `json:"from" bson:"from"`
		Gas              string `json:"gas" bson:"gas"`
		GasPrice         string `json:"gasPrice" bson:"gasPrice"`
		GasUsed          string `json:"gasUsed" bson:"gasUsed"`
		Hash             string `json:"hash" bson:"hash"`
		TxreceiptStatus  string `json:"txreceipt_status" bson:"txreceipt_status"`
		Input            string `json:"input" bson:"input"`
		Nonce            string `json:"nonce" bson:"nonce"`
		To               string `json:"to" bson:"to"`
		Type             string `json:"type" bson:"type"`
		TransactionIndex string `json:"transactionIndex" bson:"transactionIndex"`
		Value            string `json:"value" bson:"value"`
	}

	TxDecodedLogParams struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		Type  string `json:"type"`
	}

	TxDecodedLog struct {
		Address      string        `json:"address"`
		DecodedEvent *DecodedEvent `json:"decoded_event"`
	}

	DecodedEvent struct {
		Label  string                `json:"label"`
		Params []*TxDecodedLogParams `json:"params"`
	}

	TransactionDecodedReceipt struct {
		Status         string          `json:"receipt_status"`
		BlockTimestamp string          `json:"block_timestamp"`
		Logs           []*TxDecodedLog `json:"logs"`
		Message        string          `json:"message"`
	}

	BlockData struct {
		Number       string         `json:"number"`
		Hash         string         `json:"hash"`
		Nonce        string         `json:"nonce"`
		Timestamp    string         `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}

	FromToken struct {
		Value   string `json:"value" bson:"value"`
		Address string `bson:"address"`
		Token   *Token `json:"token" bson:"-"`
	}

	ToToken struct {
		Value   string `json:"value" bson:"value"`
		Address string `bson:"address"`
		Token   *Token `json:"token" bson:"-"`
	}

	SwapOrder struct {
		ID              primitive.ObjectID   `bson:"_id,omitempty"`
		ProviderKey     string               `bson:"provider_key,omitempty"`
		ChainID         uint64               `bson:"chain_id"`
		TransactionFee  string               `bson:"transaction_fee"`
		ContractAddress string               `bson:"contract_address"`
		FromAddress     string               `bson:"from_address"`
		DestReceiver    string               `bson:"dest_receiver"`
		ReceiptStatus   enum.SwapOrderStatus `bson:"receipt_status"`
		FromToken       *FromToken           `bson:"from_token,omitempty"`
		ToToken         *ToToken             `bson:"to_token,omitempty"`
		NativeToken     *Token               `bson:"-"`
		MinReturnAmount string               `bson:"min_return_amount"`
		Transaction     *Transaction         `bson:"tx,omitempty"`
		OrderAt         *time.Time           `bson:"order_at,omitempty"`
		BlockAt         *time.Time           `bson:"block_at,omitempty"`
		UpdatedAt       time.Time            `bson:"updated_at"`
		CreatedAt       time.Time            `bson:"created_at"`
		Provider        *SwapProvider        `bson:"-"`
	}
)

func (u *SwapOrder) BeforeCreate(ctx context.Context) error {
	u.CreatedAt = time.Now()
	return u.BeforeUpdate(ctx)
}

func (u *SwapOrder) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	u.FromAddress = strings.ToLower(u.FromAddress)
	u.DestReceiver = strings.ToLower(u.DestReceiver)
	u.ContractAddress = strings.ToLower(u.ContractAddress)
	if u.Transaction != nil {
		u.Transaction.Hash = strings.ToLower(u.Transaction.Hash)
	}
	if u.FromToken != nil {
		u.FromToken.Address = strings.ToLower(u.FromToken.Address)
	}
	if u.ToToken != nil {
		u.ToToken.Address = strings.ToLower(u.ToToken.Address)
	}
	return nil
}

func CreateSwapOrder(ctx context.Context, data *SwapOrder) (*SwapOrder, error) {

	if err := data.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	res, err := db.DB().Collection("swaporders").InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	data.ID = res.InsertedID.(primitive.ObjectID)
	return data, err
}

func ListSwapOrder(ctx context.Context, params IAdminParams) ([]*SwapOrder, error) {
	res := make([]*SwapOrder, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, preloadSwapOrder(ctx, res...)
}

func PageSwapOrder(ctx context.Context, params IAdminPageParams) ([]*SwapOrder, pagination.Pagination, error) {
	out := make([]*SwapOrder, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	pageParams := params.GetPageParams()
	paginator := pagination.NewTraditionalPaginator(pageParams.PageNum, pageParams.PageSize, chain)
	page, err := paginator.Paginate(&out)
	if err != nil {
		return nil, nil, err
	}
	return out, page, preloadSwapOrder(ctx, out...)
}

func FindSwapOrder(ctx context.Context, params IAdminParams) (*SwapOrder, error) {
	res := &SwapOrder{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.First(res).Error
	if err != nil {
		return nil, err
	}
	return res, preloadSwapOrder(ctx, res)
}

func preloadSwapOrder(ctx context.Context, lists ...*SwapOrder) error {
	//token
	tokenKeys := make([]string, 0)
	for _, v := range lists {
		if v.FromToken != nil {
			tokenKeys = append(tokenKeys, getTokenKey(v.ChainID, v.FromToken.Address))
		}
		if v.ToToken != nil {
			tokenKeys = append(tokenKeys, getTokenKey(v.ChainID, v.ToToken.Address))
		}
		tokenKeys = append(tokenKeys, getTokenKey(v.ChainID, NativeTokenAddress))
	}
	tokenList, err := FindSwapTokenByKeys(ctx, tokenKeys...)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	tokenMap := make(map[string]*Token)
	for _, v := range tokenList {
		token := &Token{
			Name:     v.Name,
			Address:  v.Address,
			LogoURI:  v.Logo,
			Symbol:   v.Symbol,
			Decimals: v.Decimals,
		}
		tokenMap[v.Key] = token
	}
	//provider
	providerKeys := make([]string, 0)
	for _, v := range lists {
		providerKeys = append(providerKeys, v.ProviderKey)
	}
	providerList, err := FindSwapProivderByKeys(ctx, providerKeys...)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	providerMap := make(map[string]*SwapProvider)
	for _, v := range providerList {
		providerMap[v.Key] = v
	}
	for _, v := range lists {
		if v.FromToken != nil {
			v.FromToken.Token = tokenMap[getTokenKey(v.ChainID, v.FromToken.Address)]
		}
		if v.ToToken != nil {
			v.ToToken.Token = tokenMap[getTokenKey(v.ChainID, v.ToToken.Address)]
			if v.ToToken.Value == "" {
				v.ToToken.Value = v.MinReturnAmount
			}
		}
		v.NativeToken = tokenMap[getTokenKey(v.ChainID, NativeTokenAddress)]
		v.Provider = providerMap[v.ProviderKey]
		if v.Transaction != nil {
			gasPrice, _ := new(big.Int).SetString(v.Transaction.GasPrice, 10)
			gasUsed, _ := new(big.Int).SetString(v.Transaction.GasUsed, 10)
			if gasPrice != nil && gasUsed != nil {
				v.TransactionFee = new(big.Int).Mul(gasUsed, gasPrice).String()
			}

		}
	}
	return nil
}
