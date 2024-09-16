package pkg

import (
	"chat/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	db *redis.Client
}

type KafkaStore struct {
	producer *kafka.Producer
}

type GroupStore struct {
	gs *redis.Client
}

func NewKafkaStore() (*KafkaStore, error) {
	kafkaBrokerHost := os.Getenv("KAFKA_BROKER_HOST")
	if kafkaBrokerHost == "" {
		kafkaBrokerHost = "localhost"
	}
	kafkaBrokerPort := os.Getenv("KAFKA_BROKER_PORT")
	if kafkaBrokerPort == "" {
		kafkaBrokerPort = "9092"
	}
	kafkaProducerClientId := os.Getenv("KAFKA_PRODUCER_CLIENT_ID")
	if kafkaProducerClientId == "" {
		kafkaProducerClientId = "chat-message-producer"
	}

	k, err := kafka.NewProducer(&kafka.ConfigMap{
		// "bootstrap.servers": "192.168.0.101:9092",
		"bootstrap.servers": fmt.Sprintf("%s:%s", kafkaBrokerHost, kafkaBrokerPort),
		"client.id":         kafkaProducerClientId,
		"acks":              "all"},
	)

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return nil, err
	}

	return &KafkaStore{
		producer: k,
	}, nil
}

func (ks *KafkaStore) Publish(topic string, mes types.Message) {
	fmt.Println("passed through kafka")
	messageJSON, err := json.Marshal(mes)
	if err != nil {
		log.Printf("Failed to serialize message: %v", err)
		return
	}

	err = ks.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(messageJSON)},
		nil, // delivery channel
	)

	if err != nil {
		log.Println("failed to publish message to broker", err)
		return
	}
}

func NewRedisStore() (*RedisStore, error) {
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

	return &RedisStore{
		db: rdb,
	}, nil

}

func (s *RedisStore) PubRedis(ctx context.Context, mes types.Message) {
	messageJSON, err := json.Marshal(mes)
	if err != nil {
		log.Printf("Failed to serialize message: %v", err)
		return
	}

	s.db.Publish(ctx, "msg", messageJSON)
}

// this funcition is incomplete and data like delete user, newUser needs to be handled
func (s *RedisStore) SubRedis(gpStore *GroupStore) {
	subMsg := s.db.Subscribe(context.Background(), "msg")
	msgch := subMsg.Channel()

	subGp := s.db.Subscribe(context.Background(), "groups")
	gpch := subGp.Channel()

	go func() {
		for msg := range msgch {
			var mes types.Message
			err := json.Unmarshal([]byte(msg.Payload), &mes)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			if mes.MessageType == types.GroupMessage {
				users := gpStore.getUserIdsInGroup(mes.Reciever)
				for _, user := range users {
					if UsernameToItsWebsocket[user] != nil {
						UsernameToItsWebsocket[user].WriteJSON(mes)
					}
				}
				continue
			}

			if UsernameToItsWebsocket[mes.Reciever] != nil {
				UsernameToItsWebsocket[mes.Reciever].WriteJSON(mes)
			}
		}
	}()

	go func() {
		for msg := range gpch {
			var user struct {
				// gorm.Model
				MessageType string
				GroupId     int
				GroupName   string
				GroupMember string
			}

			err := json.Unmarshal([]byte(msg.Payload), &user)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			groupUsers := gpStore.getUserIdsInGroup(user.GroupName) // we shoudl switch to using group Id later
			for _, user := range groupUsers {
				if UsernameToItsWebsocket[user] != nil {
					UsernameToItsWebsocket[user].WriteJSON(types.Message{
						MessageType: types.AddUserTOGroup,
					})
				}
			}

		}
	}()
}

func NewGroupStore() (*GroupStore, error) {

	redisGCHandlerHost := os.Getenv("REDIS_GROUP_CATALOGUE_HANDLER_HOST")
	if redisGCHandlerHost == "" {
		redisGCHandlerHost = "localhost"
	}
	redisGCHandlerPort := os.Getenv("REDIS_GROUP_CATALOGUE_HANDLER_PORT")
	if redisGCHandlerPort == "" {
		redisGCHandlerPort = "6379"
	}
	redisGCHandlerPassword := os.Getenv("REDIS_GROUP_CATALOGUE_HANDLER_PASSWORD")
	if redisGCHandlerPassword == "" {
		redisGCHandlerPassword = ""
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:             "localhost:6389",
		Password:         "",
		DB:               0,
		DisableIndentity: true, // Disable set-info on connect
	})

	return &GroupStore{
		gs: rdb,
	}, nil
}

func (gs *GroupStore) getUserIdsInGroup(groupId string) []string {
	data, err := gs.gs.LRange(context.Background(), groupId, 0, -1).Result()

	if err != nil {
		log.Fatalf("Error retrieving data from Redis Groups: %v", err)
	}

	for _, value := range data {
		fmt.Println(value)
	}
	return data
}
