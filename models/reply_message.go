package models

type ReplyMessage struct {
	Message
	reply_message_id int
	reply_text       string
	reply_sender     User
}
