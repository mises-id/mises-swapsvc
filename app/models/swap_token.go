package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mises-id/mises-swapsvc/app/models/enum"
	"github.com/mises-id/mises-swapsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Token struct {
		Address  string `json:"address" bson:"address"`
		Decimals int    `json:"decimals" bson:"decimals"`
		LogoURI  string `json:"logo_uri" bson:"logo_uri"`
		Name     string `json:"name" bson:"name"`
		Price    string `json:"price" bson:"price"`
		Symbol   string `json:"symbol" bson:"symbol"`
	}
	SwapToken struct {
		ID        primitive.ObjectID `bson:"_id,omitempty"`
		ChainID   uint64             `bson:"chain_id" json:"chain_id"`
		Status    enum.StatusType    `bson:"status" bson:"status"`
		Token     `bson:"inline"`
		Key       string    `bson:"key" json:"key"`
		UpdatedAt time.Time `bson:"updated_at"`
		CreatedAt time.Time `bson:"created_at"`
	}
)

func (u *SwapToken) BeforeCreate(ctx context.Context) error {
	u.CreatedAt = time.Now()
	u.Key = getTokenKey(u.ChainID, u.Address)
	return u.BeforeUpdate(ctx)
}

func getTokenKey(chainID uint64, address string) string {
	return fmt.Sprintf("%d&%s", chainID, strings.ToLower(address))
}

func (u *SwapToken) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	return nil
}

func CreateSwapToken(ctx context.Context, data *SwapToken) (*SwapToken, error) {

	if err := data.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	res, err := db.DB().Collection("swaptokens").InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	data.ID = res.InsertedID.(primitive.ObjectID)
	return data, err
}

func CreateSwapTokenMany(ctx context.Context, data []*SwapToken) error {
	if len(data) == 0 {
		return nil
	}
	var in []interface{}
	for _, v := range data {
		v.BeforeCreate(ctx)
		v.Status = 1
		in = append(in, v)
	}
	_, err := db.DB().Collection("swaptokens").InsertMany(ctx, in)

	return err
}

func FindSwapTokenByKeys(ctx context.Context, keys ...string) ([]*SwapToken, error) {
	res := make([]*SwapToken, 0)
	err := db.ODM(ctx).Find(&res, bson.M{"key": bson.M{"$in": keys}}).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ListSwapToken(ctx context.Context, params IAdminParams) ([]*SwapToken, error) {
	res := make([]*SwapToken, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, preloadSwapToken(ctx, res...)
}

func preloadSwapToken(ctx context.Context, lists ...*SwapToken) error {

	return nil
}
