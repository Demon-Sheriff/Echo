package models

type Room struct {
	room_id      int
	room_key     RoomKey
	room_name    string
	room_limit   int
	active_users map[int]User
}
