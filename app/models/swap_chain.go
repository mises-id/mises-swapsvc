package models

import (
	"context"
	"time"

	"github.com/mises-id/mises-swapsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	SwapChain struct {
		ID              primitive.ObjectID `bson:"_id,omitempty"`
		ChainID         uint64             `bson:"chain_id" json:"chain_id"`
		Name            string             `bson:"name" json:"name"`
		MoralisKey      string             `bson:"moralis_key" json:"moralis_key"`
		RPCAddress      string             `bson:"rpc_address" json:"rpc_address"`
		SyncBlockNumber uint64             `bson:"sync_block_number"`
		ScanApi         string             `bson:"scan_api" json:"scan_api"`
		ScanApiKey      string             `bson:"scan_api_key" json:"scan_api_key"`
		Status          int                `bson:"status"`
		SyncInterval    int                `bson:"sync_interval" json:"sync_interval"` //ms
		Tokens          []*Token           `bson:"-"`
		UpdatedAt       time.Time          `bson:"updated_at"`
		CreatedAt       time.Time          `bson:"created_at"`
	}
)

func (u *SwapChain) BeforeCreate(ctx context.Context) error {
	u.CreatedAt = time.Now()
	return u.BeforeUpdate(ctx)
}

func (u *SwapChain) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	return nil
}

func CreateSwapChain(ctx context.Context, data *SwapChain) (*SwapChain, error) {

	if err := data.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	res, err := db.DB().Collection("swapchains").InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	data.ID = res.InsertedID.(primitive.ObjectID)
	return data, err
}

func CreateSwapChainMany(ctx context.Context, data []*SwapChain) error {
	if len(data) == 0 {
		return nil
	}
	var in []interface{}
	for _, v := range data {
		v.BeforeCreate(ctx)
		v.Status = 1
		in = append(in, v)
	}
	_, err := db.DB().Collection("swapchains").InsertMany(ctx, in)

	return err
}

func ChangeSwapChainBlockNumber(ctx context.Context, id primitive.ObjectID, number uint64) error {
	update := bson.M{}
	update["sync_block_number"] = number
	_, err := db.DB().Collection("swapchains").UpdateOne(ctx, &bson.M{
		"_id": id,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}

func ListSwapChain(ctx context.Context, params IAdminParams) ([]*SwapChain, error) {
	res := make([]*SwapChain, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, preloadSwapChain(ctx, res...)
}

func FindSwapChainByChainIDs(ctx context.Context, ids ...uint64) ([]*SwapChain, error) {
	res := make([]*SwapChain, 0)
	err := db.ODM(ctx).Find(&res, bson.M{"chain_id": bson.M{"$in": ids}}).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func preloadSwapChain(ctx context.Context, lists ...*SwapChain) error {

	return nil
}
