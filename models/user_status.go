package models

type UserSTATUS int

const (
	Idle UserSTATUS = iota
	Client
	Server
	Both
)
