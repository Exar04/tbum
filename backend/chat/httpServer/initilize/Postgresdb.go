package initilize

import (
	"fmt"
	"httpChatServer/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresDB *gorm.DB

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
	PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error while initilizing DB", err)
	}
}

func SyncDb() {
	PostgresDB.AutoMigrate(&models.Messages{})
	PostgresDB.AutoMigrate(&models.Groups{})
}
