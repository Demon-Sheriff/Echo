package models

import "time"

type Message struct {
	message_id int
	sender     User
	text       string
	time_stamp time.Time
}
