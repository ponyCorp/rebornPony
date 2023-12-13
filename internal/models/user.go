package models

type User struct {
	UserId     int64
	ChatId     int64
	Role       string
	Reputation int
	Level      int
}
