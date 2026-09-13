package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	GRPC     GRPCConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Env      string `env:"APP_ENV" envDefault:"development"`
	HTTPPort string `env:"HTTP_PORT" envDefault:":8080"`
}

type GRPCConfig struct {
	Port           string        `env:"GRPC_PORT" envDefault:":50051"`
	AuthServiceURL string        `env:"AUTH_GRPC_ADDR" envDefault:"localhost:50051"` // для gateway-клиента
	Timeout        time.Duration `env:"GRPC_TIMEOUT" envDefault:"5s"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	DBName   string `env:"POSTGRES_DB" envDefault:"gateforge"`
	SSLMode  string `env:"POSTGRES_SSL" envDefault:"disable"`

	MaxConns    int32         `env:"POSTGRES_MAX_CONNS" envDefault:"10"`
	MinConns    int32         `env:"POSTGRES_MIN_CONNS" envDefault:"2"`
	MaxConnIdle time.Duration `env:"POSTGRES_MAX_CONN_IDLE" envDefault:"15m"`
	MaxConnLife time.Duration `env:"POSTGRES_MAX_CONN_LIFE" envDefault:"1h"`
}

type RedisConfig struct {
	Addr         string        `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	Password     string        `env:"REDIS_PASSWORD,required"`
	DB           int           `env:"REDIS_DB" envDefault:"0"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"3s"`
}

type JWTConfig struct {
	Secret        string        `env:"JWT_SECRET,required"`
	AccessExpiry  time.Duration `env:"JWT_ACCESS_EXPIRY" envDefault:"15m"`
	RefreshExpiry time.Duration `env:"JWT_REFRESH_EXPIRY" envDefault:"720h"` // 30 дней
}

func (p *PostgresConfig) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.User, p.Password),
		Host:   fmt.Sprintf("%s:%s", p.Host, p.Port),
		Path:   p.DBName,
	}

	q := u.Query()
	q.Set("sslmode", p.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config.Load: %w", err)
	}

	return &cfg, nil
}
