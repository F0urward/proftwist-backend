package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Mongo    MongoConfig    `yaml:"mongo"`
	Redis    RedisConfig    `yaml:"redis"`
	Service  ServiceConfig  `yaml:"service"`
	Auth     AuthConfig     `yaml:"auth"`
}

type AuthConfig struct {
	Jwt JwtConfig `yaml:"jwt"`
}

type JwtConfig struct {
	Secret string          `yaml:"secret"`
	Expire time.Duration   `yaml:"expire"`
	Cookie JwtCookieConfig `yaml:"cookie"`
}

type JwtCookieConfig struct {
	Name     string `yaml:"name"`
	MaxAge   int    `yaml:"maxAge"`
	Secure   bool   `yaml:"secure"`
	HttpOnly bool   `yaml:"httpOnly"`
}

type ServiceConfig struct {
	Host string     `yaml:"host"`
	HTTP HTTPConfig `yaml:"http"`
	CORS CORSConfig `yaml:"cors"`
}

type HTTPConfig struct {
	Port string `yaml:"port"`
}

type CORSConfig struct {
	AllowOrigin      string `yaml:"allowOrigin"`
	AllowMethods     string `yaml:"allowMethods"`
	AllowHeaders     string `yaml:"allowHeaders"`
	ExposeHeaders    string `yaml:"exposeHeaders"`
	AllowCredentials bool   `yaml:"allowCredentials"`
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

type MongoConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Driver   string `yaml:"driver"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`

	PoolSize     int           `yaml:"poolSize"`
	DialTimeout  time.Duration `yaml:"dialTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	PoolTimeout  time.Duration `yaml:"poolTimeout"`
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

		"mongo.host":     "MONGO_HOST",
		"mongo.port":     "MONGO_PORT",
		"mongo.dbname":   "MONGO_DB",
		"mongo.user":     "MONGO_USERNAME",
		"mongo.password": "MONGO_PASSWORD",

		"redis.host":     "REDIS_HOST",
		"redis.port":     "REDIS_PORT",
		"redis.db":       "REDIS_DB",
		"redis.password": "REDIS_PASSWORD",

		"auth.jwt.secret": "JWT_SECRET",
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
