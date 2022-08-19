package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/diyliv/weather/config"
	"github.com/diyliv/weather/internal/models"
	"github.com/diyliv/weather/internal/weather"
	"github.com/diyliv/weather/pkg/utils"
	"go.uber.org/zap"
)

type Handlers struct {
	logger    *zap.Logger
	weatherUC weather.PostgresUsecase
	cfg       *config.Config
}

func NewHandler(logger *zap.Logger, weatherUC weather.PostgresUsecase, cfg *config.Config) *Handlers {
	return &Handlers{logger: logger, weatherUC: weatherUC, cfg: cfg}
}

func (h *Handlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds models.Credentials

		if r.Method != "POST" {
			h.writeResponse(http.StatusBadRequest, w, "Method must be POST not "+r.Method)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			h.writeResponse(http.StatusBadRequest, w, "Invalid args")
			return
		}

		if err := h.weatherUC.Register(models.Credentials{
			Login:    creds.Login,
			Password: utils.HashPass([]byte(creds.Password)),
		}); err != nil {
			h.writeResponse(http.StatusBadRequest, w, "User with this login already exists")
			return
		}

		h.writeResponse(http.StatusOK, w, "Account was successfully created!")
	}
}

func (h *Handlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds models.Credentials

		if r.Method != "POST" {
			h.writeResponse(http.StatusBadRequest, w, "Method must be POST not "+r.Method)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			h.writeResponse(http.StatusBadRequest, w, "Invalid args")
			return
		}

		hashPass, err := h.weatherUC.Login(creds)
		if err != nil {
			h.writeResponse(http.StatusBadRequest, w, "This user does not exists. Go to the /register")
			h.logger.Error("Error while selecting creds from db: " + err.Error())
			return
		}

		if !utils.ComparePass(hashPass, []byte(creds.Password)) {
			h.writeResponse(http.StatusUnauthorized, w, "Invalid password")
			return
		}

		jwtKey := []byte(h.cfg.JWTKey.Key)

		expTime := time.Now().Add(24 * time.Hour).UTC()

		claims := &Claims{
			Username: creds.Login,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			h.logger.Error("Error while creating JWT " + err.Error())
			h.writeResponse(http.StatusInternalServerError, w, "Something went wrong :(")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    h.cfg.JWTKey.TokenName,
			Value:   tokenString,
			Expires: expTime,
		})

		if h.GetCookie(w, r, h.cfg.JWTKey.TokenName) {
			h.writeResponse(http.StatusOK, w, "Already authorized")
			return
		} else {
			h.writeResponse(http.StatusOK, w, "Successfully authorized :)")
			return
		}
	}
}

func (h *Handlers) GetForecast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			h.writeResponse(http.StatusBadRequest, w, "Method must be POST not "+r.Method)
			return
		}

		cookie := h.GetCookie(w, r, h.cfg.JWTKey.TokenName)
		if !cookie {
			return
		}

		var userCity models.City

		if err := json.NewDecoder(r.Body).Decode(&userCity); err != nil {
			h.writeResponse(http.StatusBadRequest, w, "Invalid parameters")
			return
		}

		citySlice := strings.Split(userCity.City, ",")

		if len(citySlice) < 2 {
			job := make(chan Job, 1)
			result := make(chan Result, 1)

			h.GetCoords(citySlice[0], job)

			go h.GetWeather(job, result)

			h.writeResponse(http.StatusOK, w, <-result)

		} else {
			job := make(chan Job, len(citySlice))
			result := make(chan Result, len(citySlice))

			for _, val := range citySlice {
				h.GetCoords(val, job)
			}

			val := <-job

			for w := 1; w <= len(citySlice); w++ {
				go h.GetWeather(job, result)
			}

			for i := 1; i <= len(citySlice); i++ {
				job <- Job{Lat: val.Lat, Lon: val.Lon}
			}

			close(job)

			for a := 1; a <= len(citySlice); a++ {
				h.writeResponse(http.StatusOK, w, <-result)
			}
			close(result)
		}

	}
}

func (h *Handlers) GetCookie(w http.ResponseWriter, r *http.Request, cookieName string) bool {
	c, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			h.writeResponse(http.StatusUnauthorized, w, "You must login into you acc")
			return false
		}
	}

	tokenString := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.JWTKey.Key), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			h.writeResponse(http.StatusUnauthorized, w, "Invalid creds")
			return false
		}
		h.writeResponse(http.StatusBadRequest, w, "Invalid args")
		return false
	}

	if !tkn.Valid {
		h.writeResponse(http.StatusUnauthorized, w, "Go to /login")
		return false
	}

	return true
}

func (h *Handlers) writeResponse(code int, w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		h.logger.Error("Error while encoding data: " + err.Error())
	}
}
