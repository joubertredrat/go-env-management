package pkg

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

const (
	ENV_FILE               = ".env"
	ENV_DEV                = "dev"
	ENV_STAGING            = "staging"
	ENV_PROD               = "prod"
	DB_SSLMODE_DISABLE     = "disable"
	DB_SSLMODE_ALLOW       = "allow"
	DB_SSLMODE_PREFER      = "prefer"
	DB_SSLMODE_REQUIRE     = "require"
	DB_SSLMODE_VERIFY_CA   = "verify-ca"
	DB_SSLMODE_VERIFY_FULL = "verify-full"
)

type DbSSLMode string

type Config struct {
	ApiEnv           string    `env:"API_ENV,required"`
	ApiHost          string    `env:"API_HOST,required"`
	ApiPort          string    `env:"API_PORT,required"`
	DatabaseHost     string    `env:"DATABASE_HOST,required"`
	DatabasePort     string    `env:"DATABASE_PORT,required"`
	DatabaseUsername string    `env:"DATABASE_USERNAME,required"`
	DatabasePassword string    `env:"DATABASE_PASSWORD,required"`
	DatabaseDBName   string    `env:"DATABASE_DBNAME,required"`
	DatabaseSSLMode  DbSSLMode `env:"DATABASE_SSL_MODE,required"`
}

func GetConfig() (Config, error) {
	if err := loadEnvFile(); err != nil {
		return Config{}, err
	}

	customParsers := map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(DbSSLMode("")): func(v string) (interface{}, error) {
			for _, sslmode := range dbSSLModes() {
				if sslmode == string(v) {
					return DbSSLMode(v), nil
				}
			}

			return nil, errors.New(fmt.Sprintf(
				`Invalid environment variable "DATABASE_SSL_MODE" %s, available options are: %s`,
				v,
				strings.Join(dbSSLModes(), ", "),
			))
		},
	}

	config := Config{}
	if err := env.ParseWithFuncs(&config, customParsers); err != nil {
		return Config{}, err
	}

	if !slices.Contains(envModes(), config.ApiEnv) {
		return Config{}, errors.New(fmt.Sprintf(
			`Invalid environment variable "API_ENV" %s, available options are: %s`,
			config.ApiEnv,
			strings.Join(envModes(), ", "),
		))
	}

	return config, nil
}

func loadEnvFile() error {
	if _, err := os.Stat(ENV_FILE); os.IsNotExist(err) {
		return nil
	}

	return godotenv.Load(ENV_FILE)
}

func envModes() []string {
	return []string{ENV_DEV, ENV_STAGING, ENV_PROD}
}

func dbSSLModes() []string {
	return []string{
		DB_SSLMODE_DISABLE,
		DB_SSLMODE_ALLOW,
		DB_SSLMODE_PREFER,
		DB_SSLMODE_REQUIRE,
		DB_SSLMODE_VERIFY_CA,
		DB_SSLMODE_VERIFY_FULL,
	}
}
