package models

import (
	"context"
	"strings"
	"time"

	"github.com/mises-id/mises-swapsvc/app/models/enum"
	"github.com/mises-id/mises-swapsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	SwapContract struct {
		ID              primitive.ObjectID `bson:"_id,omitempty"`
		ChainID         uint64             `bson:"chain_id" json:"chain_id"`
		Address         string             `bson:"address" json:"address"`
		ProviderKey     string             `bson:"provider_key" json:"provider_key"`
		Name            string             `bson:"name" json:"name`
		ABI             string             `bson:"abi" json:"abi`
		SyncBlockNumber uint64             `bson:"sync_block_number"`
		SyncInterval    int                `bson:"sync_interval" json:"sync_interval"` //ms
		Status          enum.StatusType    `bson:"status"`                             //1 启用 2 关闭
		Chain           *SwapChain         `bson:"-"`
		UpdatedAt       time.Time          `bson:"updated_at"`
		CreatedAt       time.Time          `bson:"created_at"`
	}
)

func (u *SwapContract) BeforeCreate(ctx context.Context) error {
	u.CreatedAt = time.Now()
	return u.BeforeUpdate(ctx)
}

func (u *SwapContract) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	u.Address = strings.ToLower(u.Address)
	return nil
}

func CreateSwapContract(ctx context.Context, data *SwapContract) (*SwapContract, error) {

	if err := data.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	res, err := db.DB().Collection("swapcontracts").InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	data.ID = res.InsertedID.(primitive.ObjectID)
	return data, err
}

func CreateSwapContractMany(ctx context.Context, data []*SwapContract) error {
	if len(data) == 0 {
		return nil
	}
	var in []interface{}
	for _, v := range data {
		v.BeforeCreate(ctx)
		v.Status = 2
		in = append(in, v)
	}
	ordered := false
	opts := &options.InsertManyOptions{
		Ordered: &ordered,
	}
	_, err := db.DB().Collection("swapcontracts").InsertMany(ctx, in, opts)

	return err
}

func ListSwapContract(ctx context.Context, params IAdminParams) ([]*SwapContract, error) {
	res := make([]*SwapContract, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, preloadSwapContract(ctx, res...)
}

func ChangeSwapContractBlockNumber(ctx context.Context, id primitive.ObjectID, number uint64) error {
	update := bson.M{}
	update["sync_block_number"] = number
	_, err := db.DB().Collection("swapcontracts").UpdateOne(ctx, &bson.M{
		"_id": id,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}

func FindSwapContract(ctx context.Context, params IAdminParams) (*SwapContract, error) {
	res := &SwapContract{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.First(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func preloadSwapContract(ctx context.Context, lists ...*SwapContract) error {
	chainIDs := make([]uint64, 0)
	for _, v := range lists {
		chainIDs = append(chainIDs, v.ChainID)
	}
	chainList, err := FindSwapChainByChainIDs(ctx, chainIDs...)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	dataMap := make(map[uint64]*SwapChain)
	for _, v := range chainList {
		dataMap[v.ChainID] = v
	}
	for _, v := range lists {
		v.Chain = dataMap[v.ChainID]
	}
	return nil
}
