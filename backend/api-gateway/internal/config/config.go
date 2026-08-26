package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port                   string
	UserServiceURL         string
	RestaurantServiceURL   string
	OrderServiceURL        string
	PaymentServiceURL      string
	NotificationServiceURL string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("[config] .env not found, using environment variables")
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:                   port,
		UserServiceURL:         viper.GetString("USER_SERVICE_URL"),
		RestaurantServiceURL:   viper.GetString("RESTAURANT_SERVICE_URL"),
		OrderServiceURL:        viper.GetString("ORDER_SERVICE_URL"),
		PaymentServiceURL:      viper.GetString("PAYMENT_SERVICE_URL"),
		NotificationServiceURL: viper.GetString("NOTIFICATION_SERVICE_URL"),
	}
}
