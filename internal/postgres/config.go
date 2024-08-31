package postgresConfig

import (
	"fmt"

	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

type PsqlInfo struct {
	Info string
}

func NewPostgresConfig(options ...func(*PostgresConfig)) *PostgresConfig {
	config := &PostgresConfig{}
	for _, option := range options {
		option(config)
	}
	return config
}

func WithHostPostgres(host string) func(*PostgresConfig) {
	return func(c *PostgresConfig) {
		c.Host = host
	}
}

func WithPortPostgres(port string) func(*PostgresConfig) {
	return func(c *PostgresConfig) {
		c.Port = port
	}
}

func WithUserPostgres(user string) func(*PostgresConfig) {
	return func(c *PostgresConfig) {
		c.User = user
	}
}

func WithPasswordPostgres(password string) func(*PostgresConfig) {
	return func(c *PostgresConfig) {
		c.Password = password
	}
}

func WithDbnamePostgres(dbname string) func(*PostgresConfig) {
	return func(c *PostgresConfig) {
		c.Dbname = dbname
	}
}

func DefaultPostgres() *PostgresConfig {
	return NewPostgresConfig(
		WithHostPostgres(utils.GetDotEnvVariable(constants.POSTGRES_HOST)),
		WithPortPostgres(utils.GetDotEnvVariable(constants.POSTGRES_PORT)),
		WithUserPostgres(utils.GetDotEnvVariable(constants.POSTGRES_USER)),
		WithPasswordPostgres(utils.GetDotEnvVariable(constants.POSTGRES_PASSWORD)),
		WithDbnamePostgres(utils.GetDotEnvVariable(constants.POSTGRES_NAME)),
	)
}

func GetPsqlInfo(config *PostgresConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)
}

func GetPsqlInfoDefault() *PsqlInfo {
	return &PsqlInfo{Info: GetPsqlInfo(DefaultPostgres())}
}
