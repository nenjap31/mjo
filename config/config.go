package config

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MysqlDBConfig struct {
	Driver   string `mapstructure:"Driver"`
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Username string `mapstructure:"Username"`
	Password string `mapstructure:"Password"`
	Name     string `mapstructure:"Name"`
}

type AppConfig struct {
	ProdMode        bool          `mapstructure:"ProdMode"`
	MysqlDb MysqlDBConfig `mapstructure:"MysqlDb"`
	SecretKey    string `mapstructure:"SecretKey"`
	JwtExpired                int `mapstructure:"JwtExpired"`
	Name         string `mapstructure:"Name"`
	RedisHost            string `mapstructure:"RedisHost"`
	RedisPassword        string `mapstructure:"RedisPassword"`
	RedisPort            string `mapstructure:"RedisPort"`
	RedisDb              int    `mapstructure:"RedisDb"`
}

func SetConfig() *AppConfig {
	var appConfig AppConfig
	viper.BindEnv("ProdMode", "PROD_MODE")
	viper.BindEnv("MysqlDb.Driver", "MASTER_DB_DRIVER")
	viper.BindEnv("MysqlDb.Host", "MASTER_DB_HOST")
	viper.BindEnv("MysqlDb.Port", "MASTER_DB_PORT")
	viper.BindEnv("MysqlDb.Username", "MASTER_DB_USERNAME")
	viper.BindEnv("MysqlDb.Password", "MASTER_DB_PASSWORD")
	viper.BindEnv("MysqlDb.Name", "MASTER_DB_NAME")
	viper.BindEnv("RedisHost", "REDIS_HOST")
	viper.BindEnv("RedisPassword", "REDIS_PASSWORD")
	viper.BindEnv("RedisPort", "REDIS_PORT")
	viper.BindEnv("RedisDb", "REDIS_DB")
	viper.BindEnv("SecretKey", "APP_SECRET_KEY")
	viper.BindEnv("JwtExpired", "JWT_EXPIRED")

	err := viper.Unmarshal(&appConfig)
	if err != nil {
		log.Info(err)
	}
	return &appConfig
}

var appConfig *AppConfig
var lock = &sync.Mutex{}

func GetConfig() *AppConfig {
	if appConfig != nil {
		return appConfig
	}

	lock.Lock()
	defer lock.Unlock()
	if appConfig != nil {
		return appConfig
	}

	appConfig = SetConfig()

	return appConfig
}