package fx_utils

import (
	"context"
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/routes"
	"github.com/nihal-ramaswamy/GoChat/internal/server"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newServerEngine(
	lc fx.Lifecycle,
	rdb_auth *redis.Client,
	config *server.Config,
	log *zap.Logger,
	pdb *sql.DB,
	ctx context.Context,
	amqpConfig *amqpConfig.AmqpConfig,
	upgrader *websocket.Upgrader,
	websocketMap *dto.WebsocketConnectionMap,
) *gin.Engine {
	gin.SetMode(config.GinMode)

	server := gin.Default()
	server.Use(cors.New(config.Cors))
	server.Use(gin.Recovery())

	routes.NewRoutes(server, pdb, rdb_auth, ctx, log, amqpConfig, upgrader, websocketMap)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting server on port", zap.String("port", config.Port))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping server")
			defer func() {
				err := log.Sync()
				if nil != err {
					log.Error(err.Error())
				}
			}()

			return nil
		},
	})

	return server
}

var serverModule = fx.Module(
	"serverModule",
	fx.Provide(
		fx.Annotate(
			newServerEngine,
			fx.ParamTags(``, `name:"rdb_auth"`),
		),
	),
)
