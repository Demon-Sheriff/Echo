package models

type User struct {
	user_id   int
	user_name string
	rooms     []Room
	status    UserSTATUS
}
