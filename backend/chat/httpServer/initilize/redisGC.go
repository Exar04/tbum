package initilize

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisGCDB *RedisGCStore

type RedisGCStore struct {
	db *redis.Client
}

func NewRedisGCStore() (*RedisWSStore, error) {
	redisWSHandlerHost := os.Getenv("REDIS_WEBSOCKET_HANDLER_HOST")
	if redisWSHandlerHost == "" {
		redisWSHandlerHost = "localhost"
	}
	redisWSHandlerPort := os.Getenv("REDIS_WEBSOCKET_HANDLER_PORT")
	if redisWSHandlerPort == "" {
		redisWSHandlerPort = "6389"
	}
	redisWSHandlerPassword := os.Getenv("REDIS_WEBSOCKET_HANDLER_PASSWORD")
	if redisWSHandlerPassword == "" {
		redisWSHandlerPassword = ""
	}
	// redisWSHandlerUser := os.Getenv("REDIS_WEBSOCKET_HANDLER_USER")
	// if redisWSHandlerHost == ""{
	// 	redisWSHandlerHost = "localhost"
	// }

	rdb := redis.NewClient(&redis.Options{
		Addr:             fmt.Sprintf("%s:%s", redisWSHandlerHost, redisWSHandlerPort),
		Password:         redisWSHandlerPassword,
		DB:               0,
		DisableIndentity: true, // Disable set-info on connect
	})

	return &RedisWSStore{
		db: rdb,
	}, nil

}

func (s *RedisGCStore) SetGroupUser(groupId, userId string) {
	err := s.db.RPush(context.Background(), groupId, userId).Err()

	if err != nil {
		log.Println("Error while adding user to existing group in redis GC", err)
	}
}
