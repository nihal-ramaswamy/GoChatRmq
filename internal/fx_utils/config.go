package fx_utils

import (
	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	postgresConfig "github.com/nihal-ramaswamy/GoChat/internal/postgres"
	redis_config "github.com/nihal-ramaswamy/GoChat/internal/redis"
	"github.com/nihal-ramaswamy/GoChat/internal/server"
	"go.uber.org/fx"
)

var ConfigModule = fx.Module(
	"Config",
	fx.Provide(server.Default),
	fx.Provide(postgresConfig.GetPsqlInfoDefault),
	fx.Provide(
		fx.Annotate(
			redis_config.DefaultRedisAuthConfig,
			fx.ResultTags(`name:"auth_rdb_config"`),
		),
	),
	fx.Provide(amqpConfig.DefaultAmqpConfig),
)
