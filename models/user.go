package models

type User struct {
	User_id   int
	User_name string
	Rooms     []Room
	Status    UserSTATUS
}
