package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type messages struct {
	gorm.Model
	Sender   string
	Reciever string
	Content  string
}

type apimsg struct {
	MessageType string `json:"messageType"`
	Sender      string `json:"sender"`
	Reciever    string `json:"reciever"`
	Content     string `json:"content"`
}

func init() {
	ConnectToDB()
	SyncDb()
}

func main() {
	kafkaPort := os.Getenv("KAFKA_BROKER_PORT")
	if kafkaPort == "" {
		kafkaPort = "9092"
	}
	kafkaHost := os.Getenv("KAFKA_BROKER_HOST")
	if kafkaHost == "" {
		kafkaHost = "localhost"
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", kafkaHost, kafkaPort),
		"group.id":          "foo",
		"auto.offset.reset": "smallest"})

	if err != nil {
		log.Fatal("Error while creating consumer", err)
	}

	err = consumer.SubscribeTopics([]string{"msg"}, nil)
	if err != nil {
		log.Fatal("Error while subscribing to topic\n", err)
	}

	run := true
	for run == true {

		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			var msg apimsg
			json.Unmarshal(e.Value, &msg)
			if msg.MessageType == "userMessage" || msg.MessageType == "groupMessage" {
				DB.Create(&messages{
					Sender:   msg.Sender,
					Reciever: msg.Reciever,
					Content:  msg.Content,
				})
			}

		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		}
	}

	consumer.Close()
}

var DB *gorm.DB

func ConnectToDB() {
	var err error

	chatPostgresPort := os.Getenv("CHAT_POSTGRES_PORT")
	if chatPostgresPort == "" {
		chatPostgresPort = "5442"
	}
	chatPostgresHost := os.Getenv("CHAT_POSTGRES_HOST")
	if chatPostgresHost == "" {
		chatPostgresHost = "localhost"
	}
	chatPostgresUser := os.Getenv("CHAT_POSTGRES_USER")
	if chatPostgresUser == "" {
		chatPostgresUser = "chat-db"
	}
	chatPostgresPassword := os.Getenv("CHAT_POSTGRES_PASSWORD")
	if chatPostgresPassword == "" {
		chatPostgresPassword = "chat-db"
	}
	chatPostgresDBName := os.Getenv("CHAT_POSTGRES_DB_NAME")
	if chatPostgresDBName == "" {
		chatPostgresDBName = "chat-db"
	}

	sslmode := "disable"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s ", //TimeZone=Asia/Shanghai
		chatPostgresHost, chatPostgresUser, chatPostgresPassword, chatPostgresDBName, chatPostgresPort, sslmode)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Error while initilizing chat DB")
	}
}

func SyncDb() {
	DB.AutoMigrate(messages{})
}
