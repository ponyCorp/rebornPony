package chat

import (
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

type Chat struct {
	driver *sqllib.ISql
}
type ChatScheme struct {
	ID string `gorm:"column:id;primary_key"`
	//chat info
	ChatName       string   `gorm:"column:chat_name"`
	ChatID         int64    `gorm:"column:chat_id"`
	WelcomeMessage string   `gorm:"column:welcome_message"`
	RulesMessage   string   `gorm:"column:rules_message"`
	KnownUsers     []string `gorm:"column:known_users;type:text[]" `
	// service
	DeleteServiceMessage        bool  `gorm:"column:delete_service_message"`
	DeleteServiceMessageTimeout int64 `gorm:"column:delete_service_message_timeout"`
	IgnoreEvents                bool  `gorm:"column:ignore_events"`
	StormMode                   bool  `gorm:"column:storm_mode"`
	//forgiveness
	ForgivenessAlert   bool  `gorm:"column:forgiveness_alert"`
	ForgivinessTimeout int64 `gorm:"column:forgiveness_timeout"`
	//tree
	ParentChat    string   `gorm:"column:parent_chat"`
	ChildrenChats []string `gorm:"column:children_chats;type:text[]" `
	EventsChannel int64    `gorm:"column:events_channel"`
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
	if result := u.driver.Driver.Where("chat_id = ?", ChatID).First(&Chat); result.Error != nil {

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
		WelcomeMessage: Chat.WelcomeMessage,
		RulesMessage:   Chat.RulesMessage,
	}, nil

}
