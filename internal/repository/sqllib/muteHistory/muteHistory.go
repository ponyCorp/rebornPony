package mutehistory

import "github.com/ponyCorp/rebornPony/internal/repository/sqllib"

type MuteHistory struct {
	driver *sqllib.ISql
}
type MuteHistoryScheme struct {
	ID      int64  `gorm:"column:id"`
	UserID  int64  `gorm:"column:user_id"`
	ChatID  int64  `gorm:"column:chat_id"`
	AdminId int64  `gorm:"column:admin_id"`
	Date    int64  `gorm:"column:date"`
	Reason  string `gorm:"column:reason"`
}

func Init(driver *sqllib.ISql) (*MuteHistory, error) {
	driver.Driver.AutoMigrate(&MuteHistoryScheme{})
	return &MuteHistory{
		driver: driver,
	}, nil
}
