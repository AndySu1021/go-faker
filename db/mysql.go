package db

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"go-faker/logger"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const BatchSize = 500

type MySQL struct {
	client *sql.DB
}

type SchemaInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default interface{}
	Extra   string
}

func (m MySQL) ParseSchema(name string) ([]string, error) {
	rows, err := m.client.Query("DESCRIBE " + name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	fields := make([]string, 0)
	for rows.Next() {
		var i SchemaInfo
		if err = rows.Scan(
			&i.Field,
			&i.Type,
			&i.Null,
			&i.Key,
			&i.Default,
			&i.Extra,
		); err != nil {
			return nil, err
		}

		if strings.Contains(i.Extra, "auto_increment") {
			continue
		}
		fields = append(fields, i.Field)
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

func (m MySQL) Create(model Model, count int, args []string) (int64, error) {
	defer func(t time.Time) {
		logger.Logger.Infow("create record log",
			"latency", fmt.Sprintf("%d ms", time.Since(t).Milliseconds()),
		)
	}(time.Now())

	def := model.Definition()
	keys := getKeys(def)
	argsMap := parseArgs(args)
	total := int64(0)
	times := math.Ceil(float64(count) / float64(BatchSize))
	wg := &sync.WaitGroup{}
	wg.Add(int(times))
	for i := 0; i < int(times); i++ {
		tmpCount := BatchSize
		if count < BatchSize {
			tmpCount = count
		}
		go func(wg *sync.WaitGroup, count int) {
			defer wg.Done()

			builder := squirrel.Insert(model.TableName()).Columns(keys...)
			for j := 0; j < count; j++ {
				def = model.Definition()
				values := getValues(def, keys, argsMap)
				builder = builder.Values(values...)
			}
			result, err := builder.RunWith(m.client).Exec()
			if err != nil {
				logger.Logger.Errorf("create record error: %s", err)
				return
			}

			// get affected rows
			affected, err := result.RowsAffected()
			if err != nil {
				logger.Logger.Errorf("get fake data affected count error: %s", err)
				return
			}

			atomic.AddInt64(&total, affected)
		}(wg, tmpCount)

		count -= BatchSize
	}
	wg.Wait()

	return total, nil
}

func parseArgs(args []string) map[string]interface{} {
	tmp := make(map[string]interface{})
	for i := 0; i < len(args); i++ {
		arr := strings.Split(args[i], "=")
		tmp[arr[0][2:]] = arr[1]
	}
	return tmp
}

func getKeys(definition map[string]interface{}) []string {
	keys := make([]string, 0)
	for k := range definition {
		keys = append(keys, k)
	}
	return keys
}

func getValues(definition map[string]interface{}, keys []string, argsMap map[string]interface{}) []interface{} {
	for k, v := range argsMap {
		if _, ok := definition[k]; ok {
			definition[k] = v
		}
	}
	values := make([]interface{}, 0)
	for i := 0; i < len(keys); i++ {
		values = append(values, definition[keys[i]])
	}
	return values
}

func (m MySQL) Truncate(model Model) error {
	str := fmt.Sprintf("TRUNCATE %s", model.TableName())
	if _, err := m.client.Exec(str); err != nil {
		return err
	}

	return nil
}

type MySQLConfig struct {
	Host               string `mapstructure:"host"`
	Port               int    `mapstructure:"port"`
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	DBName             string `mapstructure:"db_name"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	MaxLifetimeSec     int    `mapstructure:"max_lifetime_sec"`
}

func InitMySQL(cfg MySQLConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&multiStatements=true&parseTime=true", cfg.Username, cfg.Password, cfg.Host+":"+strconv.Itoa(cfg.Port), cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Logger.Errorf("fail to open connection, err: %+v", err)
		return err
	}

	if err = db.Ping(); err != nil {
		logger.Logger.Errorf("fail to ping mysql, err: %+v", err)
		return err
	}

	logger.Logger.Infof("ping mysql success")

	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxLifetimeSec) * time.Minute)

	DB = &MySQL{client: db}
	return nil
}
