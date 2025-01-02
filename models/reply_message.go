package models

type ReplyMessage struct {
	Message
	Reply_message_id int
	Reply_text       string
	Reply_sender     User
}
