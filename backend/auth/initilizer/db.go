package initilizer

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	// host := "localhost"
	host := os.Getenv("DB_HOST")
	user := "auth-db"
	password := "auth-db"
	dbName := "auth-db"
	port := "5432"
	sslmode := "disable"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s ", //TimeZone=Asia/Shanghai
		host, user, password, dbName, port, sslmode)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Error while initilizing DB")
	}
}
