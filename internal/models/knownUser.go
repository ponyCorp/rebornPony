package models

type KnownUser struct {
	UserId       int64
	ChatId       int64
	IsBot        bool
	Username     string
	FirstName    string
	LastName     string
	LanguageCode string
	FirstJoined  int64
}
