package models

type Room struct {
	Room_id      int
	Room_key     RoomKey
	Room_name    string
	Room_limit   int
	Active_users map[int]User
	Room_server Server
}
