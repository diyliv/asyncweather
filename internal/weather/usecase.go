package weather

import "github.com/diyliv/weather/internal/models"

type PostgresUsecase interface {
	Register(models.Credentials) error
	Login(models.Credentials) (string, error)
}
