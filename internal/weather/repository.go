package weather

import "github.com/diyliv/weather/internal/models"

type PostgresRepository interface {
	Register(models.Credentials) error
	Login(models.Credentials) (string, error)
}
