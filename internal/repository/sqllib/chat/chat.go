package chat

import (
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

type Chat struct {
	driver *sqllib.ISql
}
type ChatScheme struct {
	ID             string   `gorm:"column:id;primary_key"`
	ChildrenChats  []string `gorm:"column:children_chats;type:text[]" `
	ParentChat     string   `gorm:"column:parent_chat"`
	ChatID         string   `gorm:"column:chat_id"`
	DisableRead    bool     `gorm:"column:disable_read"`
	WelcomeMessage string   `gorm:"column:welcome_message"`
	RulesMessage   string   `gorm:"column:rules_message"`
}

func Init(driver *sqllib.ISql) (*Chat, error) {
	driver.Driver.AutoMigrate(&ChatScheme{})
	return &Chat{
		driver: driver,
	}, nil
}

// GetChatByID
func (u *Chat) GetChatByID(ChatID string) (models.Chat, error) {
	var Chat ChatScheme
	if result := u.driver.Driver.Where("Chat_id = ?", ChatID).First(&Chat); result.Error != nil {

		// if result.Error.Error() == "record not found" {
		// 	u.driver.Driver.Where("Chat_id = ?", 0).First(&Chat)
		// }
		return models.Chat{}, result.Error
	}

	return models.Chat{
		ID:             Chat.ID,
		ChildrenChats:  Chat.ChildrenChats,
		ParentChat:     Chat.ParentChat,
		ChatID:         Chat.ChatID,
		DisableRead:    Chat.DisableRead,
		WelcomeMessage: Chat.WelcomeMessage,
		RulesMessage:   Chat.RulesMessage,
	}, nil

}
