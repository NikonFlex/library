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

func New() (*Config, error) {
	cfg := &Config{}

	if err := cfg.readGrpcPort(); err != nil {
		return nil, err
	}
	if err := cfg.readGrpcGatewayPort(); err != nil {
		return nil, err
	}
	if err := cfg.readPGHost(); err != nil {
		return nil, err
	}
	if err := cfg.readPGUser(); err != nil {
		return nil, err
	}
	if err := cfg.readPGPort(); err != nil {
		return nil, err
	}
	if err := cfg.readPGPassword(); err != nil {
		return nil, err
	}
	if err := cfg.readPGDB(); err != nil {
		return nil, err
	}
	if err := cfg.readPGMaxConn(); err != nil {
		return nil, err
	}

	cfg.PG.URL = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		url.PathEscape(cfg.PG.User),
		url.PathEscape(cfg.PG.Password),
		net.JoinHostPort(cfg.PG.Host, cfg.PG.Port),
		cfg.PG.DB,
	)

	return cfg, nil
}

func (config *Config) readGrpcPort() error {
	var ok bool
	config.GRPC.Port, ok = os.LookupEnv("GRPC_PORT")
	if !ok || config.GRPC.Port == "" {
		return ErrGRPCPortNotSet
	}
	return nil
}

func (config *Config) readGrpcGatewayPort() error {
	var ok bool
	config.GRPC.GatewayPort, ok = os.LookupEnv("GRPC_GATEWAY_PORT")
	if !ok || config.GRPC.GatewayPort == "" {
		return ErrGRPCGatewayPortNotSet
	}
	return nil
}

func (config *Config) readPGHost() error {
	var ok bool
	config.PG.Host, ok = os.LookupEnv("POSTGRES_HOST")
	if !ok || config.PG.Host == "" {
		return ErrPostgresHostNotSet
	}
	return nil
}

func (config *Config) readPGPort() error {
	var ok bool
	config.PG.Port, ok = os.LookupEnv("POSTGRES_PORT")
	if !ok || config.PG.Port == "" {
		return ErrPostgresPortNotSet
	}
	return nil
}

func (config *Config) readPGUser() error {
	var ok bool
	config.PG.User, ok = os.LookupEnv("POSTGRES_USER")
	if !ok || config.PG.User == "" {
		return ErrPostgresUserNotSet
	}
	return nil
}

func (config *Config) readPGDB() error {
	var ok bool
	config.PG.DB, ok = os.LookupEnv("POSTGRES_DB")
	if !ok || config.PG.DB == "" {
		return ErrPostgresDBNotSet
	}
	return nil
}

func (config *Config) readPGPassword() error {
	var ok bool
	config.PG.Password, ok = os.LookupEnv("POSTGRES_PASSWORD")
	if !ok || config.PG.Password == "" {
		return ErrPostgresPasswordNotSet
	}
	return nil
}

func (config *Config) readPGMaxConn() error {
	var ok bool
	config.PG.MaxConn, ok = os.LookupEnv("POSTGRES_MAX_CONN")
	if !ok || config.PG.MaxConn == "" {
		return ErrPostgresMaxConnNotSet
	}
	return nil
}
