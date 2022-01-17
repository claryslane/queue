package mongo

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConfig struct {
	URI      string
	Database string
	User     string
	Pass     string
	Ctx      context.Context
}

var defaultConfig = MongoConfig{
	URI:      os.Getenv("$MONGO_URI"),
	Ctx:      context.Background(),
	Database: "queue",
}

func Open(config ...MongoConfig) (*mongo.Database, error) {
	conf := defaultConfig
	if config != nil {
		conf = config[0]
	}

	var client *mongo.Client
	var err error

	client, err = mongo.Connect(conf.Ctx,
		options.Client().
			ApplyURI(conf.URI),
	)
	if err != nil {
		return nil, err
	}

	return client.Database(conf.Database), client.Ping(conf.Ctx, readpref.Primary())
}
