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
	ChatId int64 `gorm:"column:chat_id;index"`
	UserId int64 `gorm:"column:user_id;index"`
	Level  int   `gorm:"column:level;index"`
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
	if err := c.driver.Driver.Where("chat_id = ? AND user_id = ?", chatId, userId).First(&admin).Error; err != nil {
		return models.Admin{}, false, err
	}
	return models.Admin{
		UserId: admin.UserId,
		ChatId: admin.ChatId,
		Level:  admintypes.AdminType(admin.Level),
	}, true, nil
}
