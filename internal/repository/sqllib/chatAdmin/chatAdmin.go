package chatadmin

import (
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/models/admintypes"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

type ChatAdmin struct {
	driver *sqllib.ISql
}
type ChatAdminScheme struct {
	ChatId int64  `gorm:"column:chat_id;index"`
	UserId int64  `gorm:"column:user_id;index"`
	Level  int    `gorm:"column:level;index"`
	Role   string `gorm:"column:role"`
}

func Init(driver *sqllib.ISql) (ChatAdmin, error) {
	driver.Driver.AutoMigrate(&ChatAdminScheme{})
	return ChatAdmin{
		driver: driver,
	}, nil
}

// GetAdminInfo(chatId int64, userId int64) (models.Admin, bool, error)
func (c ChatAdmin) GetAdminInfo(chatId int64, userId int64) (models.Admin, bool, error) {
	var admin ChatAdminScheme
	if err := c.driver.Driver.Model(&ChatAdminScheme{}).Where("chat_id = ? AND user_id = ?", chatId, userId).First(&admin).Error; err != nil {
		return models.Admin{}, false, err
	}
	return models.Admin{
		UserId: admin.UserId,
		ChatId: admin.ChatId,
		//	Level:  admintypes.AdminType(admin.Level),
	}, true, nil
}

// SetOwner(chatId int64, userId int64) error
func (c ChatAdmin) SetOwner(chatId int64, userId int64) error {
	//update or create where chat_id = ?
	err := c.driver.Driver.Model(&ChatAdminScheme{}).Where("chat_id = ?", chatId).Assign(ChatAdminScheme{
		UserId: userId,
		Role:   "owner",
		//	Level:  admintypes.Owner.Number(),
	}).FirstOrCreate(&ChatAdminScheme{
		ChatId: chatId,
		UserId: userId,
		Role:   "owner",
		//	Level:  admintypes.Owner.Number(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// SetAdminLevel(chatId int64, userId int64, level int) error
func (c ChatAdmin) SetAdminLevel(chatId int64, userId int64, level admintypes.AdminType) error {

	err := c.driver.Driver.Model(&ChatAdminScheme{}).Where("chat_id = ? AND user_id = ?", chatId, userId).Assign(ChatAdminScheme{
		//	Level: level.Number(),
	}).FirstOrCreate(&ChatAdminScheme{
		ChatId: chatId,
		UserId: userId,
		//Level:  level.Number(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}
func (c ChatAdmin) GetUserRole(chatId int64, userId int64) (admintypes.AdminType, error) {

	var admin ChatAdminScheme
	if err := c.driver.Driver.Model(&ChatAdminScheme{}).Where("chat_id = ? AND user_id = ?", chatId, userId).First(&admin).Error; err != nil {
		if err.Error() == "record not found" {
			return admintypes.User, nil
		}
		return admintypes.User, err
	}

	switch admin.Role {
	case "owner":
		return admintypes.Owner, nil
	case "moderator":
		return admintypes.Moderator, nil
	case "user":
		return admintypes.User, nil
	default:
		return admintypes.User, nil
	}

}
