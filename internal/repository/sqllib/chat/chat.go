package chat

import (
	"fmt"

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
func (u *Chat) GetChatByChatID(ChatID int64) (models.Chat, error) {
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

// GetDisabledChats() []models.Chat chats with ignore events true
func (u *Chat) GetDisabledChats() []models.Chat {
	var chats []ChatScheme
	if err := u.driver.Driver.Where("ignore_events = ?", true).Find(&chats).Error; err != nil {
		return []models.Chat{}
	}
	var res []models.Chat
	for _, chat := range chats {
		res = append(res, models.Chat{
			ID:             chat.ID,
			ChildrenChats:  chat.ChildrenChats,
			ParentChat:     chat.ParentChat,
			ChatID:         chat.ChatID,
			WelcomeMessage: chat.WelcomeMessage,
			RulesMessage:   chat.RulesMessage,
		})
	}
	return res
}

// GetChats
func (u *Chat) GetChats() []models.Chat {

	var chats []ChatScheme
	if err := u.driver.Driver.Find(&chats).Error; err != nil {
		return []models.Chat{}
	}
	var res []models.Chat
	for _, chat := range chats {
		res = append(res, models.Chat{
			ID:             chat.ID,
			ChildrenChats:  chat.ChildrenChats,
			ParentChat:     chat.ParentChat,
			ChatID:         chat.ChatID,
			WelcomeMessage: chat.WelcomeMessage,
			RulesMessage:   chat.RulesMessage,
		})
	}
	return res
}

// GetChildrenChats
func (u *Chat) GetChildrenChats(parentChatID string) []models.Chat {

	var chats []ChatScheme
	if err := u.driver.Driver.Where("parent_chat = ?", parentChatID).Find(&chats).Error; err != nil {
		return []models.Chat{}
	}
	var res []models.Chat
	for _, chat := range chats {
		res = append(res, models.Chat{
			ID:             chat.ID,
			ChildrenChats:  chat.ChildrenChats,
			ParentChat:     chat.ParentChat,
			ChatID:         chat.ChatID,
			WelcomeMessage: chat.WelcomeMessage,
			RulesMessage:   chat.RulesMessage,
		})
	}
	return res
}

// // SetUnregisterChat
// func (u *Chat) SetUnregisterChat(chatId int64) {

// 	u.driver.Driver.Model(&ChatScheme{}).Where("chat_id = ?", chatId).Update("delete_service_message", true)
// }

// EnableEvents
func (u *Chat) EnableEvents(chatId int64) error {

	return u.driver.Driver.Model(&ChatScheme{}).Where("chat_id = ?", chatId).Update("ignore_events", false).Error
}

// DisableEvents
func (u *Chat) DisableEvents(chatId int64) error {

	return u.driver.Driver.Model(&ChatScheme{}).Where("chat_id = ?", chatId).Update("ignore_events", true).Error
}

// UpdateChatRules
func (u *Chat) UpdateChatRules(chatId int64, rulesMessage string) error {

	return u.driver.Driver.Model(&ChatScheme{}).Where("chat_id = ?", chatId).Update("rules_message", rulesMessage).Error
}

// UpdateChatWelcome
func (u *Chat) UpdateChatWelcome(chatId int64, welcomeMessage string) error {

	return u.driver.Driver.Model(&ChatScheme{}).Where("chat_id = ?", chatId).Update("welcome_message", welcomeMessage).Error
}

// AddChildrenChat
func (u *Chat) AddChildrenChat(parentChatID string, childrenChatID string) error {

	return u.driver.Driver.Model(&ChatScheme{}).Where("chat_id = ?", parentChatID).Update("children_chats", childrenChatID).Error
}

// SetParentChat
func (u *Chat) SetParentChat(parentChatID string, childrenChatID string) error {

	return u.driver.Driver.Model(&ChatScheme{}).Where("chat_id = ?", childrenChatID).Update("parent_chat", parentChatID).Error
}

// // GetOrCreateChatByChatID(id int64) (models.Chat, error)
//
//	func (u *Chat) GetOrCreateChatByChatID(id int64) (models.Chat, error) {
//		existChat, err := u.GetChatByChatID(id)
//		if err != nil {
//			if err.Error() != "record not found" {
//				return models.Chat{}, err
//			}
//			return u.
//		}
//	}
//
// CreateChat(chat models.Chat) (models.Chat, error)
func (u *Chat) CreateChat(chat models.Chat) (models.Chat, error) {
	existChat, err := u.GetChatByChatID(chat.ChatID)
	if err != nil && err.Error() != "record not found" {
		return models.Chat{}, err
	}
	if existChat.ChatID == chat.ChatID {
		return existChat, fmt.Errorf("record already exist")
	}
	if err := u.driver.Driver.Create(&ChatScheme{
		ID:             chat.ID,
		ChildrenChats:  chat.ChildrenChats,
		ParentChat:     chat.ParentChat,
		ChatID:         chat.ChatID,
		WelcomeMessage: chat.WelcomeMessage,
		RulesMessage:   chat.RulesMessage,
	}).Error; err != nil {
		return models.Chat{}, err
	}
	return chat, nil

}
