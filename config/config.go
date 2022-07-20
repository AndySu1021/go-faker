package config

import (
	"github.com/spf13/viper"
	"go-faker/db"
	"go-faker/logger"
)

type AppConfig struct {
	Logger   logger.ZapConfig `mapstructure:"logger"`
	Database db.DBConfig      `mapstructure:"database"`
}

// InitConfig Initiate config
func InitConfig() (AppConfig, error) {
	viper.AutomaticEnv()
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	var config = AppConfig{}

	if err := viper.ReadInConfig(); err != nil {
		return AppConfig{}, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return AppConfig{}, err
	}

	return config, nil
}
