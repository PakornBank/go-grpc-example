package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort      string `mapstructure:"SERVER_PORT"`
	UserServiceAddr string `mapstructure:"USER_SERVICE_ADDR"`
	AuthServiceAddr string `mapstructure:"AUTH_SERVICE_ADDR"`
	CACertPath      string `mapstructure:"CA_CERT_PATH"`
	ClientCertPath  string `mapstructure:"CLIENT_CERT_PATH"`
	ClientKeyPath   string `mapstructure:"CLIENT_KEY_PATH"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
