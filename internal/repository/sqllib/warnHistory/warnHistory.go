package warnhistory

import "github.com/ponyCorp/rebornPony/internal/repository/sqllib"

type WarnHistoryScheme struct {
	UserId       int64  `gorm:"column:user_id"`
	ChatId       int64  `gorm:"column:chat_id"`
	AdminId      int64  `gorm:"column:admin_id"`
	Date         int64  `gorm:"column:date"`
	Reason       string `gorm:"column:reason"`
	SourceMsgId  int64  `gorm:"column:source_msg_id"`
	WarningMsgId int64  `gorm:"column:warning_msg_id"`
}
type WarnHistory struct {
	driver *sqllib.ISql
}

func Init(driver *sqllib.ISql) (*WarnHistory, error) {
	driver.Driver.AutoMigrate(&WarnHistoryScheme{})

	return &WarnHistory{
		driver: driver,
	}, nil
}
