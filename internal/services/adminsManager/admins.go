package adminsmanager

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/models/admintypes"
	"github.com/ponyCorp/rebornPony/internal/repository"
)

type AdminsManager struct {
	Bot *tgbotapi.BotAPI
	rep *repository.Repository
}

func New(bot *tgbotapi.BotAPI, rep *repository.Repository) *AdminsManager {
	return &AdminsManager{
		Bot: bot,
		rep: rep,
	}
}
func (a *AdminsManager) AddAdmin(chatID int64) error {
	return nil
}
func (a *AdminsManager) RemoveAdmin(chatID int64) error {
	return nil
}

// GetOwner(chatID int64) (*tgbotapi.User, error)
func (a *AdminsManager) GetOwner(chatID int64) (*tgbotapi.User, error) {
	admins, err := a.Bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	})
	if err != nil {
		return nil, err
	}
	for _, admin := range admins {
		if admin.IsCreator() {
			err := a.rep.ChatAdmin.SetOwner(chatID, admin.User.ID)
			if err != nil {
				return nil, err
			}
			return admin.User, nil
		}
	}
	return nil, fmt.Errorf("can't find owner")
}

// GetAdminInfo(chatId int64, userID int64) (*tgbotapi.ChatMember,bool, error)
func (a *AdminsManager) GetAdminInfo(chatId int64, userID int64) (models.Admin, bool, error) {
	user, isExist, err := a.rep.ChatAdmin.GetAdminInfo(chatId, userID)
	if err != nil {
		return models.Admin{}, false, err
	}
	if !isExist {
		return models.Admin{}, false, nil
	}

	return user, isExist, nil
}

// GetUserStatus(chatID int64, userID int64) admintypes.AdminType
func (a *AdminsManager) GetUserStatus(chatID int64, userID int64) (admintypes.AdminType, error) {

	admin, isExist, err := a.GetAdminInfo(chatID, userID)
	if err != nil {
		if err.Error() == "record not found" {
			return admintypes.User, nil
		}
		return admintypes.User, err
	}
	if !isExist {
		return admintypes.User, nil
	}
	return admin.Level, nil
}

// IsAdmin(chatID int64, userID int64) bool
func (a *AdminsManager) IsAdmin(chatID int64, userID int64) bool {

	admin, isExist, err := a.GetAdminInfo(chatID, userID)
	if err != nil {
		return false
	}
	if !isExist {
		return false
	}
	return admin.Role != "user"
}

// SetUserLevel(chatID int64, userID int64, level int) error
func (a *AdminsManager) SetAdminLevel(chatID int64, userID int64, level admintypes.AdminType) error {

	err := a.rep.ChatAdmin.SetAdminLevel(chatID, userID, level)
	if err != nil {
		return err
	}
	return nil
}
