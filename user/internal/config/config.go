package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort     string `mapstructure:"SERVER_PORT"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	CACertPath     string `mapstructure:"CA_CERT_PATH"`
	ServerCertPath string `mapstructure:"SERVER_CERT_PATH"`
	ServerKeyPath  string `mapstructure:"SERVER_KEY_PATH"`
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

// DBURL constructs and returns the database connection URL string
func (c *Config) DBURL() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}
