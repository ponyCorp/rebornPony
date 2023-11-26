package warnhistory

import (
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

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

func (u *WarnHistory) AddUserWarnHistory(warn models.WarnHistory) error {
	return u.driver.Driver.Create(&WarnHistoryScheme{
		UserId:       warn.UserId,
		ChatId:       warn.ChatId,
		AdminId:      warn.AdminId,
		Date:         warn.Date,
		Reason:       warn.Reason,
		SourceMsgId:  warn.SourceMsgId,
		WarningMsgId: warn.WarningMsgId,
	}).Error
}

func (u *WarnHistory) GetUserWarnHistory(userId int64) []models.WarnHistory {
	var warns []WarnHistoryScheme
	if err := u.driver.Driver.Where("user_id = ?", userId).Find(&warns).Error; err != nil {
		return nil
	}
	var result []models.WarnHistory
	for _, warn := range warns {
		result = append(result, models.WarnHistory{
			UserId:       warn.UserId,
			ChatId:       warn.ChatId,
			AdminId:      warn.AdminId,
			Date:         warn.Date,
			Reason:       warn.Reason,
			SourceMsgId:  warn.SourceMsgId,
			WarningMsgId: warn.WarningMsgId,
		})
	}
	return result
}

func (u *WarnHistory) GetUserWarnHistoryByChatId(userId int64, chatId int64) []models.WarnHistory {
	var warns []WarnHistoryScheme
	if err := u.driver.Driver.Where("user_id = ?", userId).Where("chat_id = ?", chatId).Find(&warns).Error; err != nil {
		return nil
	}
	var result []models.WarnHistory
	for _, warn := range warns {
		result = append(result, models.WarnHistory{
			UserId:       warn.UserId,
			ChatId:       warn.ChatId,
			AdminId:      warn.AdminId,
			Date:         warn.Date,
			Reason:       warn.Reason,
			SourceMsgId:  warn.SourceMsgId,
			WarningMsgId: warn.WarningMsgId,
		})
	}
	return result
}
