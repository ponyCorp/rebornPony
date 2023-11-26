package models

type Warn struct {
	Id           int
	UserId       int64
	ChatId       int64
	Count        int
	LastWarnTime int64
	History      []string
}
type WarnHistory struct {
	UserId       int64
	ChatId       int64
	AdminId      int64
	Date         int64
	Reason       string
	SourceMsgId  int64
	WarningMsgId int64
}
