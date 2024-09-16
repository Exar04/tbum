package initilize

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisWSDB *RedisWSStore

type RedisWSStore struct {
	db *redis.Client
}

func NewRedisWSStore() (*RedisWSStore, error) {
	redisWSHandlerHost := os.Getenv("REDIS_WEBSOCKET_HANDLER_HOST")
	if redisWSHandlerHost == "" {
		redisWSHandlerHost = "localhost"
	}
	redisWSHandlerPort := os.Getenv("REDIS_WEBSOCKET_HANDLER_PORT")
	if redisWSHandlerPort == "" {
		redisWSHandlerPort = "6379"
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

type GroupIssueMsg struct {
	MessageType string
	GroupName   string
	GroupMember string
	GroupId     int
}

func (s *RedisWSStore) PubRedis(ctx context.Context, mes GroupIssueMsg, publishType string) {
	messageJSON, err := json.Marshal(mes)
	if err != nil {
		log.Printf("Failed to serialize message: %v", err)
		return
	}

	s.db.Publish(ctx, publishType, messageJSON)
}
