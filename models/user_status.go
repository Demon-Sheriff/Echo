package models

type UserSTATUS int

const (
	Idle UserSTATUS = iota
	client
	server
	Both
)
