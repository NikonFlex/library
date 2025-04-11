package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigSuccess(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setup(t)
		cfg, err := New()
		require.NoError(t, err)
		require.Equal(t, "8081", cfg.GRPC.Port)
		require.Equal(t, "8080", cfg.GRPC.GatewayPort)
		require.Equal(t, "localhost", cfg.PG.Host)
		require.Equal(t, "5432", cfg.PG.Port)
		require.Equal(t, "go", cfg.PG.Password)
		require.Equal(t, "10", cfg.PG.MaxConn)
		require.Equal(t, "nikongo", cfg.PG.User)
		require.Equal(t, "godb", cfg.PG.DB)

	})
}

func TestConfigFailures(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		error  error
	}{
		{

			"GRPC_PORT not set",
			"GRPC_PORT",
			ErrGRPCPortNotSet,
		},
		{
			"GRPC_GATEWAY_PORT not set",
			"GRPC_GATEWAY_PORT",
			ErrGRPCGatewayPortNotSet,
		},
		{
			"POSTGRES_HOST not set",
			"POSTGRES_HOST",
			ErrPostgresHostNotSet,
		},
		{
			"POSTGRES_PORT not set",
			"POSTGRES_PORT",
			ErrPostgresPortNotSet,
		},
		{
			"POSTGRES_DB not set",
			"POSTGRES_DB",
			ErrPostgresDBNotSet,
		},
		{
			"POSTGRES_USER not set",
			"POSTGRES_USER",
			ErrPostgresUserNotSet,
		},
		{
			"POSTGRES_PASSWORD not set",
			"POSTGRES_PASSWORD",
			ErrPostgresPasswordNotSet,
		},
		{
			"POSTGRES_MAX_CONN not set",
			"POSTGRES_MAX_CONN",
			ErrPostgresMaxConnNotSet,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setup(t)
			t.Setenv(test.envVar, "")
			_, err := New()
			require.ErrorIs(t, err, test.error)
		})
	}
}

func setup(t *testing.T) {
	t.Setenv("GRPC_PORT", "8081")
	t.Setenv("GRPC_GATEWAY_PORT", "8080")
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_DB", "godb")
	t.Setenv("POSTGRES_USER", "nikongo")
	t.Setenv("POSTGRES_PASSWORD", "go")
	t.Setenv("POSTGRES_MAX_CONN", "10")
}
