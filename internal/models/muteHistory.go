package models

type MuteHistory struct {
	UserId  int64
	ChatId  int64
	AdminId int64
	Date    int64
	Reason  string
}
