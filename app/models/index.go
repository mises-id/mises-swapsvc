package models

import (
	"context"
	"time"

	"github.com/mises-id/mises-swapsvc/lib/db"
	"github.com/mises-id/mises-swapsvc/lib/db/odm"
	"github.com/mises-id/mises-swapsvc/lib/pagination"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type (
	IAdminParams interface {
		BuildAdminSearch(chain *odm.DB) *odm.DB
	}
	IAdminPageParams interface {
		BuildAdminSearch(chain *odm.DB) *odm.DB
		GetPageParams() *pagination.TraditionalParams
	}
	IAdminQuickPageParams interface {
		BuildAdminSearch(chain *odm.DB) *odm.DB
		GetQuickPageParams() *pagination.PageQuickParams
	}
)

func EnsureIndex() {
	tokenIndexName := "uniqueChainAndAddress"
	opts := options.CreateIndexes().SetMaxTime(20 * time.Second)
	trueBool := true
	_, err := db.DB().Collection("swapcontracts").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{
				Key: "chain_id", Value: bsonx.Int32(1),
			}, {
				Key: "address", Value: bsonx.Int32(1)},
			},
			Options: &options.IndexOptions{
				Unique: &trueBool,
			},
		},
	}, opts)
	if err != nil {
		logrus.Debug(err)
	}
	_, err = db.DB().Collection("swaptokens").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{
				Key: "chain_id", Value: bsonx.Int32(1),
			}, {
				Key: "address", Value: bsonx.Int32(1)},
			},
			Options: &options.IndexOptions{
				Unique: &trueBool,
				Name:   &tokenIndexName,
			},
		},
	}, opts)
	if err != nil {
		logrus.Debug(err)
	}
	_, err = db.DB().Collection("swaporders").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{
				Key: "chain_id", Value: bsonx.Int32(1),
			}, {
				Key: "tx.hash", Value: bsonx.Int32(1)},
			},
			Options: &options.IndexOptions{
				Unique: &trueBool,
			},
		},
	}, opts)
	if err != nil {
		logrus.Debug(err)
	}
	_, err = db.DB().Collection("swapchains").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{"chain_id": 1},
			Options: &options.IndexOptions{
				Unique: &trueBool,
			},
		},
	}, opts)
	if err != nil {
		logrus.Debug(err)
	}
	if err != nil {
		logrus.Debug(err)
	}
	_, err = db.DB().Collection("swapproviders").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{"key": 1},
			Options: &options.IndexOptions{
				Unique: &trueBool,
			},
		},
	}, opts)
	if err != nil {
		logrus.Debug(err)
	}
}
