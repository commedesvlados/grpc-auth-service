package models

type App struct {
	ID     int
	Name   string
	Secret string // TODO: move app secret from struct
}
