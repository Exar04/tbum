package main

import (
	"chat/pkg"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {

	port := os.Getenv("CHAT_WEBSOCKET_SERVER_PORT")

	if port == "" {
		port = "8088"
	}

	Rstore, err := pkg.NewRedisStore()
	if err != nil {
		log.Fatal(err)
	}

	Groupstore, err := pkg.NewGroupStore()
	if err != nil {
		log.Fatal(err)
	}

	Kstore, err := pkg.NewKafkaStore()
	if err != nil {
		log.Fatal(err)
	}

	api := &pkg.APIServer{
		ClientConnections: make(map[*websocket.Conn]struct{}),
		Rstore:            Rstore,
		Kstore:            Kstore,
	}

	go Rstore.SubRedis(Groupstore)

	http.HandleFunc("/websocket", api.SocketHandler)

	log.Println("Starting server on localhost:8088")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	fmt.Println("hope this works!")
}
