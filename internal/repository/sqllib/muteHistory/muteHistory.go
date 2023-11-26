package mutehistory

import (
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

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

// AddUserMuteHistory(mute models.MuteHistory) error
func (m *MuteHistory) AddUserMuteHistory(mute models.MuteHistory) error {
	return m.driver.Driver.Create(&mute).Error
}

// GetUserMuteHistory(userId int64) []models.MuteHistory
func (m *MuteHistory) GetUserMuteHistory(userId int64) ([]models.MuteHistory, error) {
	var mute []MuteHistoryScheme
	if err := m.driver.Driver.Where("user_id = ?", userId).Find(&mute).Error; err != nil {
		return []models.MuteHistory{}, err
	}
	var res []models.MuteHistory
	for _, v := range mute {
		res = append(res, models.MuteHistory{

			UserId:  v.UserID,
			ChatId:  v.ChatID,
			AdminId: v.AdminId,
			Date:    v.Date,
			Reason:  v.Reason,
		})
	}
	return res, nil
}

// GetUserMuteHistoryByChatId(userId int64, chatId int64) ([]models.MuteHistory, error)
func (m *MuteHistory) GetUserMuteHistoryByChatId(userId int64, chatId int64) ([]models.MuteHistory, error) {
	var mute []MuteHistoryScheme
	if err := m.driver.Driver.Where("user_id = ?", userId).Where("chat_id = ?", chatId).Find(&mute).Error; err != nil {
		return []models.MuteHistory{}, err
	}
	var res []models.MuteHistory
	for _, v := range mute {
		res = append(res, models.MuteHistory{
			UserId:  v.UserID,
			ChatId:  v.ChatID,
			AdminId: v.AdminId,
			Date:    v.Date,
			Reason:  v.Reason,
		})
	}
	return res, nil
}
