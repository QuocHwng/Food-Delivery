package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWT         JWTConfig
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpire  time.Duration
	RefreshExpire time.Duration
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("[config] .env not found, using environment variables")
	}

	accessMinutes := viper.GetInt("JWT_ACCESS_EXPIRE_MINUTES")
	if accessMinutes == 0 {
		accessMinutes = 15
	}
	refreshDays := viper.GetInt("JWT_REFRESH_EXPIRE_DAYS")
	if refreshDays == 0 {
		refreshDays = 7
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8081"
	}

	dbURL := viper.GetString("DATABASE_URL")
	if dbURL == "" {
		// fallback: build from individual vars
		dbURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			viper.GetString("POSTGRES_USER"),
			viper.GetString("POSTGRES_PASSWORD"),
			viper.GetString("POSTGRES_HOST"),
			viper.GetString("POSTGRES_PORT"),
			viper.GetString("POSTGRES_DB"),
		)
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWT: JWTConfig{
			AccessSecret:  viper.GetString("JWT_ACCESS_SECRET"),
			RefreshSecret: viper.GetString("JWT_REFRESH_SECRET"),
			AccessExpire:  time.Duration(accessMinutes) * time.Minute,
			RefreshExpire: time.Duration(refreshDays) * 24 * time.Hour,
		},
	}
}
