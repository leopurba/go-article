package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort          int           `mapstructure:"APP_PORT"`
	Timeout          int           `mapstructure:"TIMEOUT"`
	PostgresUser     string        `mapstructure:"POSTGRES_USER"`
	PostgresPassword string        `mapstructure:"POSTGRES_PASSWORD"`
	PostgresHost     string        `mapstructure:"POSTGRES_HOST"`
	PostgresPort     int           `mapstructure:"POSTGRES_PORT"`
	PostgresDatabase string        `mapstructure:"POSTGRES_DATABASE"`
	RedisPassword    string        `mapstructure:"REDIS_PASSWORD"`
	RedisHost        string        `mapstructure:"REDIS_HOST"`
	RedisPort        int           `mapstructure:"REDIS_PORT"`
	RedisDatabase    int           `mapstructure:"REDIS_DATABASE"`
	RedisPoolSize    int           `mapstructure:"REDIS_POOLSIZE"`
	RedisTTL         time.Duration `mapstructure:"REDIS_TTL"`
}

func load(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func Cfg() *Config {
	config, err := load(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	return &config
}
