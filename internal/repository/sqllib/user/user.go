package user

import (
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/models/roles"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	driver *sqllib.ISql
}

func Init(driver *sqllib.ISql) (*User, error) {
	driver.Driver.AutoMigrate(&UserScheme{})
	return &User{
		driver: driver,
	}, nil
}

type UserScheme struct {
	ChatId     int64  `gorm:"column:chat_id;index"`
	UserId     int64  `gorm:"column:user_id;index"`
	Role       string `gorm:"column:role"`
	Level      int    `gorm:"column:level;index"`
	Reputation int    `gorm:"column:reputation"`
}

// GetOrCreateUser(chatId int64, userId int64) (*models.User, error)
func (u *User) GetOrCreateUser(chatId int64, userId int64) (*models.User, error) {
	var user UserScheme
	if err := u.driver.Driver.Model(&UserScheme{}).Where("chat_id = ? AND user_id = ?", chatId, userId).Attrs(UserScheme{Role: "user", Reputation: 0, Level: 1}).FirstOrCreate(&user).Error; err != nil {
		return nil, err
	}
	return &models.User{
		UserId:     user.UserId,
		ChatId:     user.ChatId,
		Role:       user.Role,
		Reputation: user.Reputation,
		Level:      user.Level,
	}, nil
}

// IncreaseReputation(chatId int64, userId int64, inc int) error
func (u *User) IncreaseReputation(chatId int64, userId int64, inc int) (int, error) {
	var user UserScheme
	err := u.driver.Driver.Model(&user).Clauses(clause.Returning{
		Columns: []clause.Column{
			{Name: "reputation"},
		},
	}).Where("chat_id = ? AND user_id = ?", chatId, userId).Update("reputation", gorm.Expr("reputation + ?", inc)).Error
	return user.Reputation, err
}

// DecreaseReputation(chatId int64, userId int64, inc int) error
func (u *User) DecreaseReputation(chatId int64, userId int64, inc int) (int, error) {
	var user UserScheme
	err := u.driver.Driver.Model(&user).Clauses(clause.Returning{
		Columns: []clause.Column{
			{Name: "reputation"},
		},
	}).Where("chat_id = ? AND user_id = ?", chatId, userId).Update("reputation", gorm.Expr("reputation - ?", inc)).Error
	return user.Reputation, err
}

// SetReputation(chatId int64, userId int64, reputation int) error
func (u *User) SetReputation(chatId int64, userId int64, reputation int) (int, error) {
	var user UserScheme
	err := u.driver.Driver.Model(&user).Clauses(clause.Returning{
		Columns: []clause.Column{
			{Name: "reputation"},
		},
	}).Where("chat_id = ? AND user_id = ?", chatId, userId).Update("reputation", reputation).Error
	return user.Reputation, err

}

// GetUserReputation(chatId int64, userId int64) (int, error)
func (u *User) GetUserReputation(chatId int64, userId int64) (int, error) {
	var user UserScheme
	err := u.driver.Driver.Model(&UserScheme{}).Where("chat_id = ? AND user_id = ?", chatId, userId).Select("reputation").First(&user).Error
	return user.Reputation, err
}

// SetOwner(chatId int64, userId int64) error
func (u *User) SetRole(chatId int64, userId int64, role roles.Role) error {
	return u.driver.Driver.Model(&UserScheme{}).Where("chat_id = ? AND user_id = ?", chatId, userId).Assign(UserScheme{
		Role:  role.String(),
		Level: 1,
	}).FirstOrCreate(&UserScheme{
		ChatId: chatId,
		UserId: userId,
		Role:   role.String(),
		Level:  1,
	}).Error
}
