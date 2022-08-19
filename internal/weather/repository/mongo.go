package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/diyliv/weather/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	logDB         = "log"
	logCollection = "log"
)

type mongoRepo struct {
	logger *zap.Logger
	mongo  *mongo.Client
}

func NewMongoRepo(logger *zap.Logger, mongo *mongo.Client) *mongoRepo {
	return &mongoRepo{
		logger: logger,
		mongo:  mongo,
	}
}

func (r *mongoRepo) Write(p []byte) (n int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	collection := r.mongo.Database(logDB).Collection(logCollection)

	var logMessage models.LogMessage

	if err := json.Unmarshal(p, &logMessage); err != nil {
		r.logger.Error("Error while unmarshalling: " + err.Error())
	}

	_, err = collection.InsertOne(ctx, logMessage)
	if err != nil {
		r.logger.Error("Error while inserting data: " + err.Error())
	}

	return len(p), nil
}
