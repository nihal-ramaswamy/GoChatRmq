package fx_utils

import (
	"github.com/nihal-ramaswamy/GoChat/internal/db"
	"go.uber.org/fx"
)

var postgresModule = fx.Module(
	"PostgresService",
	fx.Provide(db.GetPostgresDbInstanceWithConfig),
)
