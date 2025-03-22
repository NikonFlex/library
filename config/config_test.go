package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig_Success(t *testing.T) {
	GRPC_PORT := "50051"
	GRPC_GATEWAY_PORT := "8080"
	POSTGRES_DB := "godb"
	POSTGRES_USER := "nikongo"
	POSTGRES_PASSWORD := "go"
	POSTGRES_PORT := "5432"
	POSTGRES_HOST := "localhost"
	POSTGRES_MAX_CONN := "10"

	os.Setenv("GRPC_PORT", GRPC_PORT)
	os.Setenv("GRPC_GATEWAY_PORT", GRPC_GATEWAY_PORT)
	os.Setenv("POSTGRES_DB", POSTGRES_DB)
	os.Setenv("POSTGRES_USER", POSTGRES_USER)
	os.Setenv("POSTGRES_PASSWORD", POSTGRES_PASSWORD)
	os.Setenv("POSTGRES_PORT", POSTGRES_PORT)
	os.Setenv("POSTGRES_HOST", POSTGRES_HOST)
	os.Setenv("POSTGRES_MAX_CONN", POSTGRES_MAX_CONN)

	defer func() {
		os.Unsetenv("GRPC_PORT")
		os.Unsetenv("GRPC_GATEWAY_PORT")
		os.Unsetenv("POSTGRES_DB")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_MAX_CONN")
	}()

	cfg, err := NewConfig()
	require.NoError(t, err)
	require.Equal(t, cfg.GRPC.Port, GRPC_PORT)
	require.Equal(t, cfg.GRPC.GatewayPort, GRPC_GATEWAY_PORT)
	require.Equal(t, cfg.PG.Port, POSTGRES_PORT)
	require.Equal(t, cfg.PG.DB, POSTGRES_DB)
	require.Equal(t, cfg.PG.Host, POSTGRES_HOST)
	require.Equal(t, cfg.PG.User, POSTGRES_USER)
	require.Equal(t, cfg.PG.Password, POSTGRES_PASSWORD)
	require.Equal(t, cfg.PG.MaxConn, POSTGRES_MAX_CONN)
}

func TestNewConfig_Failure(t *testing.T) {
	vars := []string{"GRPC_PORT", "GRPC_GATEWAY_PORT", "POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_PORT", "POSTGRES_HOST", "POSTGRES_MAX_CONN"}
	tests := []struct {
		name          string
		envVar        string
		expectedError error
	}{
		{
			name:          "GRPC_PORT not set",
			envVar:        "GRPC_PORT",
			expectedError: ErrGRPCPortNotSet,
		},
		{
			name:          "GRPC_GATEWAY_PORT not set",
			envVar:        "GRPC_GATEWAY_PORT",
			expectedError: ErrGRPCGatewayPortNotSet,
		},
		{
			name:          "POSTGRES_PORT not set",
			envVar:        "POSTGRES_PORT",
			expectedError: ErrPostgresPortNotSet,
		},
		{
			name:          "POSTGRES_DB not set",
			envVar:        "POSTGRES_DB",
			expectedError: ErrPostgresDBNotSet,
		},
		{
			name:          "POSTGRES_HOST not set",
			envVar:        "POSTGRES_HOST",
			expectedError: ErrPostgresHostNotSet,
		},
		{
			name:          "POSTGRES_USER not set",
			envVar:        "POSTGRES_USER",
			expectedError: ErrPostgresUserNotSet,
		},
		{
			name:          "POSTGRES_PASSWORD not set",
			envVar:        "POSTGRES_PASSWORD",
			expectedError: ErrPostgresPasswordNotSet,
		},
		{
			name:          "POSTGRES_MAX_CONN not set",
			envVar:        "POSTGRES_MAX_CONN",
			expectedError: ErrPostgresMaxConnNotSet,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			for _, envVar := range vars {
				os.Setenv(envVar, "random")
			}
			os.Unsetenv(test.envVar)

			_, err := NewConfig()
			require.ErrorIs(t, err, test.expectedError)
		})
	}
}
