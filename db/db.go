package db

import (
	"errors"
)

type IDB interface {
	Create(model Model, count int, args []string) (int64, error)
	Truncate(model Model) error
	ParseSchema(name string) ([]string, error)
}

type Model interface {
	TableName() string
	Definition() map[string]interface{}
}

var DB IDB

type DatabaseType string

const (
	DatabaseTypeMySQL DatabaseType = "mysql"
)

type DBConfig struct {
	Type               DatabaseType `mapstructure:"type"`
	Host               string       `mapstructure:"host"`
	Port               int          `mapstructure:"port"`
	Username           string       `mapstructure:"username"`
	Password           string       `mapstructure:"password"`
	DBName             string       `mapstructure:"db_name"`
	MaxIdleConnections int          `mapstructure:"max_idle_connections"`
	MaxOpenConnections int          `mapstructure:"max_open_connections"`
	MaxLifetimeSec     int          `mapstructure:"max_lifetime_sec"`
}

func InitDatabase(cfg DBConfig) error {
	switch cfg.Type {
	case DatabaseTypeMySQL:
		return InitMySQL(MySQLConfig{
			Host:               cfg.Host,
			Port:               cfg.Port,
			Username:           cfg.Username,
			Password:           cfg.Password,
			DBName:             cfg.DBName,
			MaxIdleConnections: cfg.MaxIdleConnections,
			MaxOpenConnections: cfg.MaxOpenConnections,
			MaxLifetimeSec:     cfg.MaxLifetimeSec,
		})
	}

	return errors.New("wrong database type")
}
