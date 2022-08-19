package server

import (
	"github.com/diyliv/weather/config"
	"github.com/diyliv/weather/internal/weather"
	handlers "github.com/diyliv/weather/internal/weather/delivery/http/v1"
	"go.uber.org/zap"
)

type server struct {
	logger    *zap.Logger
	cfg       *config.Config
	weatherUC weather.PostgresUsecase
	handlers  *handlers.Handlers
}

func NewServer(logger *zap.Logger, cfg *config.Config, weatherUC weather.PostgresUsecase, handlers *handlers.Handlers) *server {
	return &server{
		logger:    logger,
		cfg:       cfg,
		weatherUC: weatherUC,
		handlers:  handlers,
	}
}
