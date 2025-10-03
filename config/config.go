package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Service  ServiceConfig  `yaml:"service"`
}

type ServiceConfig struct {
	Host string     `yaml:"host"`
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Port string `yaml:"port"`
}

type PostgresConfig struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	Driver          string `yaml:"driver"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	ConnMaxLifeTime int    `yaml:"connMaxLifeTime"`
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime"`
}

func New() *Config {
	viper, err := newViper()
	if err != nil {
		log.Fatalf("cannot create config: %v", err)
	}

	cfg, err := parseConfig(viper)
	if err != nil {
		log.Fatalf("cannot parse config: %v", err)
	}

	return cfg
}

func newViper() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(os.Getenv("CONFIG_PATH"))
	v.SetConfigName(os.Getenv("CONFIG_FILE"))
	v.SetConfigType("yml")

	err := bindEnv(v)
	if err != nil {
		return nil, fmt.Errorf("cannot bind env variables: %v", err)
	}

	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Config search paths:", v.ConfigFileUsed())
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found")
		}

		return nil, err
	}

	return v, nil
}

func bindEnv(v *viper.Viper) error {
	envBindings := map[string]string{
		"postgres.host":     "POSTGRES_HOST",
		"postgres.port":     "POSTGRES_PORT",
		"postgres.dbname":   "POSTGRES_DB",
		"postgres.user":     "POSTGRES_USER",
		"postgres.password": "POSTGRES_PASSWORD",
	}

	for key, env := range envBindings {
		if err := v.BindEnv(key, env); err != nil {
			return err
		}
	}

	return nil
}

func parseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config

	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %v", err)
	}

	return &cfg, nil
}
