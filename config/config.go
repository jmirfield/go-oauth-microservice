package config

import (
	"github.com/spf13/viper"
)

// Config contains all of the variables required by the auth service
type Config struct {
	Port string
	DSN  string
	Key  string
}

// LoadConfig returns Config struct
func LoadConfig() *Config {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("DSN", "host=localhost port=5432 user=postgres password=password dbname=auth sslmode=disable")
	viper.SetDefault("KEY_FILE", "private.pem")

	cfg := &Config{
		Port: viper.GetString("PORT"),
		DSN:  viper.GetString("DSN"),
		Key:  viper.GetString("KEY_FILE"),
	}

	return cfg
}
