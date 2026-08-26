package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	RabbitMQURL string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("[config] .env not found, using environment variables")
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8085"
	}

	return &Config{
		Port:        port,
		DatabaseURL: viper.GetString("DATABASE_URL"),
		JWTSecret:   viper.GetString("JWT_ACCESS_SECRET"),
		RabbitMQURL: viper.GetString("RABBITMQ_URL"),
	}
}
