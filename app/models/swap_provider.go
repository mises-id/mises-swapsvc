package models

import (
	"context"
	"time"

	"github.com/mises-id/mises-swapsvc/app/models/enum"
	"github.com/mises-id/mises-swapsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	SwapProvider struct {
		ID        primitive.ObjectID `bson:"_id,omitempty"`
		Name      string             `bson:"name,omitempty"`
		Key       string             `bson:"key"`
		Logo      string             `bson:"logo,omitempty"`
		Status    enum.StatusType    `bson:"status" bson:"status"`
		UpdatedAt time.Time          `bson:"updated_at"`
		CreatedAt time.Time          `bson:"created_at"`
	}
)

func (u *SwapProvider) BeforeCreate(ctx context.Context) error {
	u.CreatedAt = time.Now()
	return u.BeforeUpdate(ctx)
}

func (u *SwapProvider) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	return nil
}

func CreateSwapProvider(ctx context.Context, data *SwapProvider) (*SwapProvider, error) {

	if err := data.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	res, err := db.DB().Collection("swapproviders").InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	data.ID = res.InsertedID.(primitive.ObjectID)
	return data, err
}

func CreateSwapProviderMany(ctx context.Context, data []*SwapProvider) error {
	if len(data) == 0 {
		return nil
	}
	var in []interface{}
	for _, v := range data {
		v.BeforeCreate(ctx)
		v.Status = 1
		in = append(in, v)
	}
	_, err := db.DB().Collection("swapproviders").InsertMany(ctx, in)

	return err
}

func ListSwapProvider(ctx context.Context, params IAdminParams) ([]*SwapProvider, error) {
	res := make([]*SwapProvider, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, preloadSwapProvider(ctx, res...)
}

func FindSwapProivderByKeys(ctx context.Context, keys ...string) ([]*SwapProvider, error) {
	res := make([]*SwapProvider, 0)
	err := db.ODM(ctx).Find(&res, bson.M{"key": bson.M{"$in": keys}}).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FindSwapProvider(ctx context.Context, params IAdminParams) (*SwapProvider, error) {
	res := &SwapProvider{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.First(res).Error
	if err != nil {
		return nil, err
	}
	return res, preloadSwapProvider(ctx, res)
}

func preloadSwapProvider(ctx context.Context, lists ...*SwapProvider) error {

	return nil
}
