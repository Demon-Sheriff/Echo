package models

type ClientSTATUS int

const (
	Idle ClientSTATUS = iota
	client
	server
	Both
)
