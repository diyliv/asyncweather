package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/diyliv/weather/internal/models"
)

type Job struct {
	Lat float64
	Lon float64
}

type Result struct {
	models.Weather
}

func (h *Handlers) GetCoords(cityName string, job chan<- Job) (float64, float64) {
	if strings.Contains(cityName, " ") {
		cityName = strings.Replace(cityName, " ", "+", -1)
	}

	var latLon []models.Coords

	resp, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s", cityName, h.cfg.OpenWeather.APIKey))
	if err != nil {
		h.logger.Error("Error while interacting with OpenWeatherAPI for getting geo cords" + err.Error())
	}

	if err := json.NewDecoder(resp.Body).Decode(&latLon); err != nil {
		panic(err)
	}

	for _, val := range latLon {
		job <- Job{Lat: val.Latitude, Lon: val.Longitude}
		return val.Latitude, val.Longitude
	}

	return 0.0, 0.0
}

func (h *Handlers) GetWeather(job <-chan Job, result chan<- Result) models.Weather {
	val := <-job

	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&lang=ru&units=metric&appid=%s",
		val.Lat,
		val.Lon,
		h.cfg.OpenWeather.APIKey))
	if err != nil {
		h.logger.Error("Error while decoding body: " + err.Error())
	}

	var weatherForecast models.Weather

	if err := json.NewDecoder(resp.Body).Decode(&weatherForecast); err != nil {
		h.logger.Error("Error while decoding body: " + err.Error())
	}

	result <- Result{weatherForecast}
	return weatherForecast
}
