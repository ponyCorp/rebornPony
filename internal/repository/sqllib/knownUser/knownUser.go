package knownuser

import (
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

type KnownUser struct {
	driver *sqllib.ISql
}
type KnownUserScheme struct {
	Id     int64 `gorm:"column:id"`
	UserId int64 `gorm:"column:user_id"`
	ChatId int64 `gorm:"column:chat_id"`

	Username     string `gorm:"column:username"`
	FirstName    string `gorm:"column:first_name"`
	LastName     string `gorm:"column:last_name"`
	LanguageCode string `gorm:"column:language_code"`
	FirstJoined  int64  `gorm:"column:first_joined"`

	IsBot bool `gorm:"column:is_bot"`
}

func Init(driver *sqllib.ISql) (*KnownUser, error) {
	driver.Driver.AutoMigrate(&KnownUserScheme{})
	return &KnownUser{
		driver: driver,
	}, nil
}

// GetKnownUser
func (u *KnownUser) GetKnownUser(userId int64, chatId int64) (models.KnownUser, error) {
	var user KnownUserScheme
	if result := u.driver.Driver.Where("user_id = ?", userId).Where("chat_id = ?", chatId).First(&user); result.Error != nil {
		return models.KnownUser{}, result.Error
	}

	return models.KnownUser{
		UserId:       user.UserId,
		ChatId:       user.ChatId,
		IsBot:        user.IsBot,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		LanguageCode: user.LanguageCode,
		FirstJoined:  user.FirstJoined,
	}, nil

}

// AddKnownUser
func (u *KnownUser) AddKnownUser(knownUser models.KnownUser) (models.KnownUser, error) {

	user := KnownUserScheme{
		UserId:       knownUser.UserId,
		ChatId:       knownUser.ChatId,
		IsBot:        knownUser.IsBot,
		Username:     knownUser.Username,
		FirstName:    knownUser.FirstName,
		LastName:     knownUser.LastName,
		LanguageCode: knownUser.LanguageCode,
		FirstJoined:  knownUser.FirstJoined,
	}
	_, err := u.GetKnownUser(knownUser.UserId, knownUser.ChatId)
	if err != nil {
		return models.KnownUser{}, err
	}
	if err := u.driver.Driver.Model(&KnownUserScheme{}).Create(&user).Error; err != nil {
		return models.KnownUser{}, err
	}
	return models.KnownUser{
		UserId:       user.UserId,
		ChatId:       user.ChatId,
		IsBot:        user.IsBot,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		LanguageCode: user.LanguageCode,
		FirstJoined:  user.FirstJoined,
	}, nil
}

// RemoveKnownUser
func (u *KnownUser) RemoveKnownUser(knownUser models.KnownUser) error {

	return u.driver.Driver.Where("user_id = ?", knownUser.UserId).Where("chat_id = ?", knownUser.ChatId).Delete(&KnownUserScheme{}).Error
}
