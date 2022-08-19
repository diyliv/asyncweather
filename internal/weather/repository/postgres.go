package repository

import (
	"database/sql"

	"github.com/diyliv/weather/internal/models"
	"go.uber.org/zap"
)

type postgresRepository struct {
	psql   *sql.DB
	logger *zap.Logger
}

func NewPostgresRepository(psql *sql.DB, logger *zap.Logger) *postgresRepository {
	return &postgresRepository{
		psql:   psql,
		logger: logger,
	}
}

func (r *postgresRepository) Register(creds models.Credentials) error {
	_, err := r.psql.Exec("INSERT INTO users (login, password) VALUES($1, $2)", creds.Login, creds.Password)
	if err != nil {
		r.logger.Error("Error while creating user: " + err.Error())
		return err
	}
	return nil
}

func (r *postgresRepository) Login(creds models.Credentials) (string, error) {
	var hashPass string

	err := r.psql.QueryRow("SELECT password FROM users WHERE login = $1", creds.Login).Scan(&hashPass)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Info("This user does not exists." + err.Error())
			return "", err
		}
		r.logger.Error("Error while getting user creds: " + err.Error())
	}

	return hashPass, nil
}
