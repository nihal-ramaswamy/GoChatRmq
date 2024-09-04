package testUtils

import (
	"context"
	"fmt"

	rdb "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func SetUpRedisForTesting(ctx context.Context) (*redis.RedisContainer, *rdb.Client, error) {
	redisContainer, err := redis.Run(ctx,
		"docker.io/redis:7",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error setting up redis container: %s", err)
	}

	uri, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting connection string: %s", err)
	}

	rdbOptions, err := rdb.ParseURL(uri)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing url: %s", err)
	}

	rdbConn := rdb.NewClient(rdbOptions)

	if err := rdbConn.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("error pinging redis: %s", err)
	}

	return redisContainer, rdbConn, nil
}

func ReadFromRedis(rdb *rdb.Client, key string) (bool, string, error) {
	exists := rdb.Exists(context.Background(), key).Val()
	if exists == 0 {
		return false, "", nil
	}
	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		return true, "", err
	}
	return true, val, nil
}
