package types

type ApiMessageType string

const (
	NewUser             ApiMessageType = "newUser"
	UserMessage         ApiMessageType = "userMessage"
	GroupMessage        ApiMessageType = "groupMessage"
	DeleteUserFromGroup ApiMessageType = "deleteUserFromGroup"
	AddUserTOGroup      ApiMessageType = "addUserTOGroup"
	DeleteMessage       ApiMessageType = "deleteMessage"
)

type Message struct {
	MessageType ApiMessageType `json:"messageType"`
	Content     string         `json:"content"`
	Sender      string         `json:"sender"`
	Reciever    string         `json:"reciever"`
}
