package usecase

import (
	"github.com/diyliv/weather/internal/models"
	"github.com/diyliv/weather/internal/weather"
)

type weatherUC struct {
	postgresRepo weather.PostgresUsecase
}

func NewWeatherUC(postgresRepo weather.PostgresRepository) *weatherUC {
	return &weatherUC{postgresRepo: postgresRepo}
}

func (u *weatherUC) Register(creds models.Credentials) error {
	return u.postgresRepo.Register(creds)
}
func (u *weatherUC) Login(creds models.Credentials) (string, error) {
	return u.postgresRepo.Login(creds)
}
