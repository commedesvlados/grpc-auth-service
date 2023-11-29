package models

type App struct {
	ID     int
	Name   string
	Secret string // TODO: mode app secret from struct
}
