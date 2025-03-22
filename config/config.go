package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
)

type (
	Config struct {
		GRPC
		PG
	}

	GRPC struct {
		Port        string `env:"GRPC_PORT"`
		GatewayPort string `env:"GRPC_GATEWAY_PORT"`
	}

	PG struct {
		URL      string
		Host     string `env:"POSTGRES_HOST"`
		Port     string `env:"POSTGRES_PORT"`
		DB       string `env:"POSTGRES_DB"`
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		MaxConn  string `env:"POSTGRES_MAX_CONN"`
	}
)

var (
	ErrGRPCPortNotSet         = errors.New("GRPC_PORT environment variable not set")
	ErrGRPCGatewayPortNotSet  = errors.New("GRPC_GATEWAY_PORT environment variable not set")
	ErrPostgresHostNotSet     = errors.New("POSTGRES_HOST environment variable not set")
	ErrPostgresPortNotSet     = errors.New("POSTGRES_PORT environment variable not set")
	ErrPostgresDBNotSet       = errors.New("POSTGRES_DB environment variable not set")
	ErrPostgresUserNotSet     = errors.New("POSTGRES_USER environment variable not set")
	ErrPostgresPasswordNotSet = errors.New("POSTGRES_PASSWORD environment variable not set")
	ErrPostgresMaxConnNotSet  = errors.New("POSTGRES_MAX_CONN environment variable not set")
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	grpcPort, err := getEnv("GRPC_PORT")
	if err != nil {
		return nil, ErrGRPCPortNotSet
	}
	cfg.GRPC.Port = grpcPort

	grpcGatewayPort, err := getEnv("GRPC_GATEWAY_PORT")
	if err != nil {
		return nil, ErrGRPCGatewayPortNotSet
	}
	cfg.GRPC.GatewayPort = grpcGatewayPort

	pgHost, err := getEnv("POSTGRES_HOST")
	if err != nil {
		return nil, ErrPostgresHostNotSet
	}
	cfg.PG.Host = pgHost

	pgPort, err := getEnv("POSTGRES_PORT")
	if err != nil {
		return nil, ErrPostgresPortNotSet
	}
	cfg.PG.Port = pgPort

	pgDB, err := getEnv("POSTGRES_DB")
	if err != nil {
		return nil, ErrPostgresDBNotSet
	}
	cfg.PG.DB = pgDB

	pgUser, err := getEnv("POSTGRES_USER")
	if err != nil {
		return nil, ErrPostgresUserNotSet
	}
	cfg.PG.User = pgUser

	pgPassword, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, ErrPostgresPasswordNotSet
	}
	cfg.PG.Password = pgPassword

	pgMaxConn, err := getEnv("POSTGRES_MAX_CONN")
	if err != nil {
		return nil, ErrPostgresMaxConnNotSet
	}
	cfg.PG.MaxConn = pgMaxConn

	cfg.PG.URL = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		url.PathEscape(cfg.PG.User),
		url.PathEscape(cfg.PG.Password),
		net.JoinHostPort(cfg.PG.Host, cfg.PG.Port),
		cfg.PG.DB,
	)

	return cfg, nil
}

func getEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return value, nil
}
