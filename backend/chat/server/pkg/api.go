package pkg

import (
	"chat/types"
	"context"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type APIServer struct {
	ClientConnections map[*websocket.Conn]struct{}
	Kstore            *KafkaStore
	Rstore            *RedisStore
}

var (
	upgrader = websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by returning true or
			// implement your logic to check the origin here
			return true
		},
	}
	// use to store the temperory data of user connected via websocket
	WebSocketUsername      = make(map[*websocket.Conn]string)
	UsernameToItsWebsocket = make(map[string]*websocket.Conn)
)

func (s *APIServer) SocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("there was a connection error : ", err)
		return
	}
	defer ws.Close()

	s.ClientConnections[ws] = struct{}{}

	for {
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			s.handleDisconnection(ws)
			break
		}
		err1 := s.handleIncomingMessage(ws, bytes)
		if err1 != nil {
			log.Print("Error handling message", err1)
		}
	}
}

func (s *APIServer) handleIncomingMessage(sender *websocket.Conn, msg []byte) error {
	var DataRecieved types.Message
	err := json.Unmarshal(msg, &DataRecieved)
	if err != nil {
		return err
	}
	fmt.Println(DataRecieved)

	switch DataRecieved.MessageType {
	case types.NewUser:
		WebSocketUsername[sender] = DataRecieved.Sender
		UsernameToItsWebsocket[DataRecieved.Sender] = sender
	case types.UserMessage, types.GroupMessage:
		s.Kstore.Publish("msg", DataRecieved)
		s.Rstore.PubRedis(context.Background(), DataRecieved)
	}
	return nil
}

// don't know if this works properly
func (s *APIServer) handleDisconnection(sender *websocket.Conn) {
	// user_id, _ := WebSocketUserId[s.webSocket] // this should be used in future and remove the below ones
	user_id, _ := WebSocketUsername[sender]
	delete(s.ClientConnections, sender)
	delete(WebSocketUsername, sender)
	delete(UsernameToItsWebsocket, user_id)
}
