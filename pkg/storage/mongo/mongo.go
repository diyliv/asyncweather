package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/diyliv/weather/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnMongo(cfg *config.Config) *mongo.Client {
	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.Mongo.Login,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.Port)).
		SetConnectTimeout(time.Duration(cfg.Mongo.ConnTimeout) * time.Second).
		SetMaxConnIdleTime(time.Duration(cfg.Mongo.MaxIdleConnTime) * time.Minute).
		SetMinPoolSize(uint64(cfg.Mongo.MinPoolSize)).
		SetMaxPoolSize(uint64(cfg.Mongo.MaxPoolSize))

	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatalf("Error while connecting to mongo DB: %v\n", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error while pinging mongo DB: %v\n", err)
	}

	return client
}
