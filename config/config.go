package config

import (
	"github.com/spf13/viper"
)

// Config contains all of the variables required by the auth service
type Config struct {
	Port           string
	DSN            string
	PrivateKeyPath string
}

// LoadConfig returns Config struct
func LoadConfig() *Config {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", 3000)
	viper.SetDefault("DSN", "host=localhost port=5432 user=postgres password=password dbname=auth sslmode=disable")
	viper.SetDefault("PRIVATE_KEY_PATH", "./certificates/private.pem")

	cfg := &Config{
		Port:           viper.GetString("PORT"),
		DSN:            viper.GetString("DSN"),
		PrivateKeyPath: viper.GetString("PRIVATE_KEY_PATH"),
	}

	return cfg
}
