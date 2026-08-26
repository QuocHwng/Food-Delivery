package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	Port            string
	DatabaseURL     string
	VnpTmnCode      string
	VnpHashSecret   string
	VnpUrl          string
	VnpReturnUrl    string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("[config] .env not found, using environment variables")
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8084"
	}

	return &Config{
		Port:          port,
		DatabaseURL:   viper.GetString("DATABASE_URL"),
		VnpTmnCode:    viper.GetString("VNPAY_TMN_CODE"),
		VnpHashSecret: viper.GetString("VNPAY_HASH_SECRET"),
		VnpUrl:        viper.GetString("VNPAY_URL"),
		VnpReturnUrl:  viper.GetString("VNPAY_RETURN_URL"),
	}
}
