package db

import (
	"context"
	"database/sql"

	postgresConfig "github.com/nihal-ramaswamy/GoChat/internal/postgres"
	redis_config "github.com/nihal-ramaswamy/GoChat/internal/redis"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GetPostgresDbInstanceWithConfig(psqlInfo *postgresConfig.PsqlInfo, log *zap.Logger) *sql.DB {
	log.Info("Connecting to database...")
	db, err := sql.Open("postgres", psqlInfo.Info)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Pinging postgres...")
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

func GetRedisDbInstanceWithConfig(redisInfo *redis_config.RedisConfig, log *zap.Logger, ctx context.Context) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisInfo.Addr,
		Password: redisInfo.Password,
		DB:       redisInfo.DB,
	})

	log.Info("Pinging redis...")
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err.Error())
	}

	return rdb
}
