package main

import (
	"github.com/diyliv/weather/config"
	"github.com/diyliv/weather/internal/server"
	handlers "github.com/diyliv/weather/internal/weather/delivery/http/v1"
	"github.com/diyliv/weather/internal/weather/repository"
	"github.com/diyliv/weather/internal/weather/usecase"
	"github.com/diyliv/weather/pkg/logger"
	"github.com/diyliv/weather/pkg/storage/mongo"
	"github.com/diyliv/weather/pkg/storage/postgres"
	"github.com/diyliv/weather/pkg/storage/redis"
)

func main() {
	cfg := config.ReadConfig()

	psql := postgres.ConnPostgres(cfg)
	mongo := mongo.ConnMongo(cfg)
	redis.ConnRedis(cfg)
	logger := logger.InitLogger(mongo)

	weatherPostgresRepo := repository.NewPostgresRepository(psql, logger)
	weatherUC := usecase.NewWeatherUC(weatherPostgresRepo)
	handler := handlers.NewHandler(logger, weatherUC, cfg)
	server := server.NewServer(logger, cfg, weatherUC, handler)

	server.StartHTTP()
}
