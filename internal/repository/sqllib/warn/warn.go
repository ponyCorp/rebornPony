package warn

import (
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
	"gorm.io/gorm"
)

type WarnScheme struct {
	Id           int      `gorm:"column:id"`
	UserId       int64    `gorm:"column:user_id"`
	ChatId       int64    `gorm:"column:chat_id"`
	Count        int      `gorm:"column:count"`
	LastWarnTime int64    `gorm:"column:last_warn_time"`
	History      []string `gorm:"foreignKey:WarnID"`
}

type Warn struct {
	driver *sqllib.ISql
}

func Init(driver *sqllib.ISql) (*Warn, error) {
	driver.Driver.AutoMigrate(&WarnScheme{})

	return &Warn{
		driver: driver,
	}, nil
}
func (u *Warn) GetWarnedUserByChatId(userId int64, chatId int64) (models.Warn, error) {
	user, err := u.getOrCreateWarnedUser(userId, chatId)
	if err != nil {
		return models.Warn{}, err
	}
	return models.Warn{
		Id:           user.Id,
		UserId:       user.UserId,
		ChatId:       user.ChatId,
		Count:        user.Count,
		LastWarnTime: user.LastWarnTime,
		History:      user.History,
	}, nil
}
func (u *Warn) getOrCreateWarnedUser(userId int64, chatId int64) (WarnScheme, error) {
	var user WarnScheme
	if result := u.driver.Driver.Where("user_id = ?", userId).Where("chat_id = ?", chatId).First(&user); result.Error != nil {
		if result.Error.Error() == "record not found" {
			u.driver.Driver.Model(&WarnScheme{}).Create(&WarnScheme{
				UserId:       userId,
				ChatId:       chatId,
				Count:        0,
				LastWarnTime: 0,
			})
		}

		return WarnScheme{}, result.Error
	}
	return user, nil
}

// IncreaseWarn
func (u *Warn) IncreaseWarn(userId int64, chatId int64) (int, error) {
	user, err := u.getOrCreateWarnedUser(userId, chatId)
	if err != nil {
		return 0, err
	}
	if err = u.driver.Driver.Model(&WarnScheme{}).Where("user_id = ?", userId).Where("chat_id = ?", chatId).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
		return 0, err
	}
	return user.Count + 1, nil
}

// DecreaseWarn
func (u *Warn) DecreaseWarn(userId int64, chatId int64) (int, error) {
	user, err := u.getOrCreateWarnedUser(userId, chatId)
	if err != nil {
		return 0, err
	}
	if user.Count == 0 {
		return 0, nil
	}
	if err := u.driver.Driver.Model(&WarnScheme{}).Where("user_id = ?", userId).Where("chat_id = ?", chatId).Update("count", gorm.Expr("count - ?", 1)).Error; err != nil {
		return 0, err
	}
	return user.Count - 1, nil
}

// ResetWarn
func (u *Warn) ResetWarn(userId int64, chatId int64) error {
	user, err := u.getOrCreateWarnedUser(userId, chatId)
	if err != nil {
		return err
	}
	if user.Count == 0 {
		return nil
	}
	return u.driver.Driver.Model(&WarnScheme{}).Where("user_id = ?", userId).Where("chat_id = ?", chatId).Update("count", 0).Error
}

// GetUserWarnsFromAllChats(userId int64) ([]models.Warn,error)
func (u *Warn) GetUserWarnsFromAllChats(userId int64) ([]models.Warn, error) {
	var warns []WarnScheme
	if err := u.driver.Driver.Where("user_id = ?", userId).Find(&warns).Error; err != nil {
		return nil, err
	}
	var result []models.Warn
	for _, warn := range warns {
		result = append(result, models.Warn{
			Id:           warn.Id,
			UserId:       warn.UserId,
			ChatId:       warn.ChatId,
			Count:        warn.Count,
			LastWarnTime: warn.LastWarnTime,
			History:      warn.History,
		})
	}
	return result, nil
}
