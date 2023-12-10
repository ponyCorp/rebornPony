package chatsensetive

import (
	"fmt"

	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/models/sensetivetypes"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

type ChatSensetive struct {
	driver *sqllib.ISql
}
type ChatScheme struct {
	ID            string `gorm:"column:id;primary_key"`
	ChatID        int64  `gorm:"column:chat_id"`
	Word          string `gorm:"column:word"`
	SensetiveType string `gorm:"column:sensetive_type"`
}

func Init(driver *sqllib.ISql) (*ChatSensetive, error) {
	driver.Driver.AutoMigrate(&ChatScheme{})
	return &ChatSensetive{
		driver: driver,
	}, nil
}

// AddChatSensetive(chatId int64, sensetiveWord string) error
func (u *ChatSensetive) AddChatSensetive(chatId int64, sensetiveWord string, sensetiveType sensetivetypes.SensetiveType) error {

	return u.driver.Driver.Model(&ChatScheme{}).Create(&ChatScheme{
		ChatID:        chatId,
		Word:          sensetiveWord,
		SensetiveType: sensetiveType.String(),
	}).Error
}

// RemoveChatSensetive(chatId int64, sensetiveWord string) error
func (u *ChatSensetive) RemoveChatSensetive(chatId int64, sensetiveWord string, sensetiveType sensetivetypes.SensetiveType) error {

	return u.driver.Driver.Model(&ChatScheme{}).Where("chat_id = ? AND word = ? AND sensetive_type = ?", chatId, sensetiveWord, sensetiveType.String()).Delete(&ChatScheme{}).Error
}

// GetChatSensetiveWords(chatId int64) models.ChatSensetiveWords
func (u *ChatSensetive) GetChatSensetiveWords(chatId int64, sensetiveType sensetivetypes.SensetiveType) models.ChatSensetiveWords {

	var sensetiveWords []ChatScheme
	err := u.driver.Driver.Where("chat_id = ? AND sensetive_type = ?", chatId, sensetiveType.String()).Find(&sensetiveWords).Error
	if err != nil {
		fmt.Printf("GetChatSensetiveWords: %v\n", err)
		return models.ChatSensetiveWords{}
	}
	var res models.ChatSensetiveWords
	for _, v := range sensetiveWords {
		res.Words = append(res.Words, v.Word)
	}
	res.ChatID = chatId
	return res
}

// GetAllChatSensetiveWords(sensetiveType sensetivetypes.SensetiveType) []models.ChatSensetiveWords
func (u *ChatSensetive) GetAllChatSensetiveWords(sensetiveType sensetivetypes.SensetiveType) map[int64][]string {
	var sensetiveWords []ChatScheme
	err := u.driver.Driver.Where("sensetive_type = ?", sensetiveType.String()).Find(&sensetiveWords).Error
	if err != nil {
		fmt.Printf("GetAllChatSensetiveWords: %v\n", err)
		return map[int64][]string{}
	}

	resMap := make(map[int64][]string)
	for _, v := range sensetiveWords {
		if resMap[v.ChatID] == nil {
			resMap[v.ChatID] = []string{v.Word}
		} else {
			resMap[v.ChatID] = append(resMap[v.ChatID], v.Word)
		}
	}

	return resMap
}
