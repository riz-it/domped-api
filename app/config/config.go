package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database Database
	Logger   Logger
	Server   Server
	Jwt      JWTConfig
	Redis    Redis
	SMTP     SMTP
	Midtrans Midtrans
}

type Server struct {
	Name    string
	Version string
	Host    string
	Port    string
}

type JWTConfig struct {
	AccessTokenKey  string
	AccessTokenExp  string
	RefreshTokenKey string
	RefreshTokenExp string
}

type Database struct {
	Host                  string
	Port                  string
	Name                  string
	User                  string
	Pass                  string
	Tz                    string
	IdleConnection        string
	MaxConnection         string
	MaxLifeTimeConnection string
}

type Logger struct {
	Level string
}

type Redis struct {
	Address string
	User    string
	Pass    string
	DB      string
}

type SMTP struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Midtrans struct {
	Key    string
	IsProd bool
}

func Get() *Config {
	err := godotenv.Load()
	if err != nil {
		// log.Fatal("error load configuration: ", err.Error())
	}

	return &Config{
		Server: Server{
			Name:    os.Getenv("APP_NAME"),
			Host:    os.Getenv("APP_HOST"),
			Version: os.Getenv("APP_VERSION"),
			Port:    os.Getenv("APP_PORT"),
		},
		Logger: Logger{
			Level: os.Getenv("LOG_LEVEL"),
		},
		Database: Database{
			Host:                  os.Getenv("DB_HOST"),
			Port:                  os.Getenv("DB_PORT"),
			User:                  os.Getenv("DB_USER"),
			Pass:                  os.Getenv("DB_PASS"),
			Name:                  os.Getenv("DB_NAME"),
			Tz:                    os.Getenv("DB_TZ"),
			IdleConnection:        os.Getenv("DB_POOL_IDLE"),
			MaxConnection:         os.Getenv("DB_POOL_MAX"),
			MaxLifeTimeConnection: os.Getenv("DB_POOL_LIFETIME"),
		},
		Jwt: JWTConfig{
			AccessTokenKey:  os.Getenv("JWT_ACCESS_KEY"),
			AccessTokenExp:  os.Getenv("JWT_ACCESS_EXP"),
			RefreshTokenKey: os.Getenv("JWT_REFRESH_KEY"),
			RefreshTokenExp: os.Getenv("JWT_REFRESH_EXP"),
		},
		Redis: Redis{
			Address: os.Getenv("REDIS_ADDR"),
			User:    os.Getenv("REDIS_USER"),
			Pass:    os.Getenv("REDIS_PASS"),
			DB:      os.Getenv("REDIS_DB"),
		},
		SMTP: SMTP{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			User:     os.Getenv("SMTP_USER"),
			Password: os.Getenv("SMTP_PASS"),
		},
		Midtrans: Midtrans{
			Key:    os.Getenv("MIDTRANS_KEY"),
			IsProd: os.Getenv("MIDTRANS_ENV") == "production",
		},
	}
}
