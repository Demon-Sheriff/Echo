package models

import "time"

type Message struct {
	Message_id int
	Sender     User
	Text       string
	Time_stamp time.Time
}
