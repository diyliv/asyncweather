package config

import "github.com/spf13/viper"

type Config struct {
	Postgres    Postgres
	Mongo       Mongo
	Redis       Redis
	HttpServer  HttpServer
	JWTKey      JWTKey
	OpenWeather OpenWeather
}

type Postgres struct {
	Host            string
	Port            string
	Login           string
	Password        string
	ConnMaxLifeTime int
	MaxOpenConn     int
	MaxIdleConn     int
}

type Mongo struct {
	Host            string
	Port            string
	Login           string
	Password        string
	ConnTimeout     int
	MaxIdleConnTime int
	MinPoolSize     int
	MaxPoolSize     int
}

type Redis struct {
	Addr        string
	Password    string
	DB          int
	MinIdleConn int
	PoolSize    int
	PoolTimeout int
}

type HttpServer struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

type JWTKey struct {
	TokenName string
	Key       string
}

type OpenWeather struct {
	APIKey string
}

func ReadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
