package fx_utils

import (
	"context"

	"github.com/nihal-ramaswamy/GoChat/internal/db"
	"go.uber.org/fx"
)

var cacheModule = fx.Module(
	"CacheService",
	fx.Provide(func() context.Context {
		return context.Background()
	}),

	fx.Provide(
		fx.Annotate(
			db.GetRedisDbInstanceWithConfig,
			fx.ParamTags(`name:"auth_rdb_config"`),
			fx.ResultTags(`name:"rdb_auth"`),
		),
	),
)
