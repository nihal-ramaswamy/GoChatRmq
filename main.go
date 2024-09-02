package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	"github.com/nihal-ramaswamy/GoChat/internal/fx_utils"
	"github.com/nihal-ramaswamy/GoChat/internal/server"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(utils.NewZapLogger),
		utils.FxLogger(),

		fx_utils.WebsocketModule,

		fx_utils.ConfigModule,
		fx_utils.MicroServicesModule,

		fx.Invoke(Invoke),
	).Run()
}

func Invoke(
	server *gin.Engine,
	config *server.Config,
	amqpConfig *amqpConfig.AmqpConfig,
	log *zap.Logger,
) {
	defer func() {
		log.Info("Closing amqp connection")
		amqpConfig.Channel.Close()
		amqpConfig.Conn.Close()
	}()

	err := server.Run(config.Port)
	if nil != err {
		log.Error(err.Error())
	}
}
