package db

import (
	"context"

	"github.com/mises-id/mises-swapsvc/config/env"
	"github.com/mises-id/mises-swapsvc/lib/db/odm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoDB     *mongo.Database
	mongoClient *mongo.Client
	mongoDBMap  map[string]*mongo.Database
	odmClient   *odm.Client
)

func SetupMongo(ctx context.Context) {
	clientOpts := options.Client().SetMaxPoolSize(30).ApplyURI(env.Envs.MongoURI)
	if env.Envs.DBUser != "" {
		clientOpts = clientOpts.SetAuth(options.Credential{
			Username: env.Envs.DBUser,
			Password: env.Envs.DBPass,
		})
	}
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	mongoClient = client
	mongoDB = client.Database(env.Envs.DBName)
	odmClient = odm.NewClient(mongoDB)
}

func DB() *mongo.Database {
	return mongoDB
}
func Database(dbname string) *mongo.Database {
	var res *mongo.Database
	res, ok := mongoDBMap[dbname]
	if !ok {
		res = mongoClient.Database(dbname)
	}
	return res
}
func ODM(ctx context.Context) *odm.DB {
	return odmClient.NewSession(ctx)
}
func NewODM(ctx context.Context, dbname string) *odm.DB {
	return odm.NewClient(Database(dbname)).NewSession(ctx)
}
