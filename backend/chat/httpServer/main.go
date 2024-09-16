package main

import (
	"context"
	"encoding/json"
	"httpChatServer/initilize"
	"httpChatServer/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// this server will handle following things -
// - /getOldMessages
// - /deleteMessage
// - /createGroup
// - /addOrRemoveUserFromGroup
// - /sync with redis-groups-db
// - manageLastSceneMessage and userLastActiveAt

func init() {
	initilize.ConnectToDB()
	initilize.SyncDb()

	var err error
	initilize.RedisWSDB, err = initilize.NewRedisWSStore()
	if err != nil {
		log.Fatal("error while connecting to redis websocket server", err)
	}
}

func main() {
	port := os.Getenv("CHAT_HTTP_SERVER_PORT")
	if port == "" {
		port = "9055"
	}
	router := mux.NewRouter()

	router.HandleFunc("/getMessages/{chatId}/{fromDate}", getMessages).Methods(http.MethodGet)
	router.HandleFunc("/deleteMessage/{messageId}", deleteMessage).Methods(http.MethodDelete)
	router.HandleFunc("/getGroupCatalogue", getGroupCatalogue).Methods(http.MethodGet)
	router.HandleFunc("/getGroupUsers/{groupId}", getGroupUsers).Methods(http.MethodGet)
	router.HandleFunc("/addUserToGroup/{groupId}", addUserInGroup).Methods(http.MethodPost)

	log.Println("ChatDB service listining on port", port)
	http.ListenAndServe(":"+port, router)
}

// imp feature
func getMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId := vars["chatId"]
	fromDate := vars["fromDate"]

	createdAtDate, err := time.Parse("2006-01-02T15:04:05Z", fromDate)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	var messages []models.Messages

	// Query messages based on chatId and createdAtDate
	if err := initilize.PostgresDB.Where("created_at > ? AND (sender = ? OR receiver = ?)", createdAtDate, chatId, chatId).Find(&messages).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the messages in JSON format
	data, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// imp feature
func getGroupCatalogue(w http.ResponseWriter, r *http.Request) {
	var groups []models.Groups

	// Fetch all groups from the database
	if err := initilize.PostgresDB.Find(&groups).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the groups in JSON format
	data, err := json.Marshal(groups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getGroupUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]

	var users []models.Groups
	err := initilize.PostgresDB.Where("group_name = ?", groupId).Find(&users).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// imp feature
func addUserInGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]

	var user struct {
		GroupName   string
		GroupMember string
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &user)
	groupIdint, err := strconv.ParseInt(groupId, 6, 12)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newUser := models.Groups{
		GroupName:   user.GroupName,
		GroupMember: user.GroupMember,
		GroupId:     int(groupIdint),
	}
	initilize.PostgresDB.Create(&newUser)

	var addUser struct {
		MessageType string
		GroupName   string
		GroupMember string
		GroupId     int
	}
	addUser.MessageType = "addUserToGroup"
	addUser.GroupName = user.GroupName
	addUser.GroupMember = user.GroupMember
	addUser.GroupId = int(groupIdint)

	initilize.RedisWSDB.PubRedis(context.Background(), addUser, "groups")
	initilize.RedisGCDB.SetGroupUser(groupId, user.GroupMember)
}

// after deleting message in database it should also send delete message request to redis sub
// so that it will flow through all the websocktes and delete message in realtime if user is online
func deleteMessage(w http.ResponseWriter, r *http.Request) {

}
func userLastActiveAt(w http.ResponseWriter, r *http.Request) {

}

func messageLastScene(w http.ResponseWriter, r *http.Request) {

}
