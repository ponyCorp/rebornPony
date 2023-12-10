package repository

import (
	"fmt"

	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/models/sensetivetypes"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
	sqllibChat "github.com/ponyCorp/rebornPony/internal/repository/sqllib/chat"
	sqllibKnownUser "github.com/ponyCorp/rebornPony/internal/repository/sqllib/knownUser"
	sqllibMuteHistory "github.com/ponyCorp/rebornPony/internal/repository/sqllib/muteHistory"

	sqllibChatsensetive "github.com/ponyCorp/rebornPony/internal/repository/sqllib/chatSensetive"
	sqllibWarn "github.com/ponyCorp/rebornPony/internal/repository/sqllib/warn"
	sqllibWarnHistory "github.com/ponyCorp/rebornPony/internal/repository/sqllib/warnHistory"
)

type KnownUser interface {
	AddKnownUser(knownUser models.KnownUser) (models.KnownUser, error)
	RemoveKnownUser(knownUser models.KnownUser) error
	GetKnownUser(userId int64, chatId int64) (models.KnownUser, error)
}
type Warn interface {
	IncreaseWarn(userId int64, chatId int64) (int, error)
	DecreaseWarn(userId int64, chatId int64) (int, error)
	ResetWarn(userId int64, chatId int64) error
	GetWarnedUserByChatId(userId int64, chatId int64) (models.Warn, error)
	GetUserWarnsFromAllChats(userId int64) ([]models.Warn, error)
}
type WarnHistory interface {
	AddUserWarnHistory(warn models.WarnHistory) error
	GetUserWarnHistory(userId int64) []models.WarnHistory
	GetUserWarnHistoryByChatId(userId int64, chatId int64) []models.WarnHistory
}
type MuteHistory interface {
	AddUserMuteHistory(mute models.MuteHistory) error
	GetUserMuteHistory(userId int64) ([]models.MuteHistory, error)
	GetUserMuteHistoryByChatId(userId int64, chatId int64) ([]models.MuteHistory, error)
}

type Chat interface {
	CreateChat(chat models.Chat) (models.Chat, error)
	GetChatByChatID(id int64) (models.Chat, error)
	GetDisabledChats() []models.Chat
	GetChats() []models.Chat
	GetChildrenChats(parentChatID string) []models.Chat
	//SetUnregisterChat(chatId int64) //deprecated
	//SetRegisterChat(chatId int64)//deprecated
	EnableEvents(chatId int64) error
	DisableEvents(chatId int64) error
	UpdateChatRules(chatId int64, rulesMessage string) error
	UpdateChatWelcome(chatId int64, welcomeMessage string) error
	AddChildrenChat(parentChatID string, childrenChatID string) error
	SetParentChat(parentChatID string, childrenChatID string) error
}
type ChatSensetive interface {
	AddChatSensetive(chatId int64, sensetiveWord string, sensetiveType sensetivetypes.SensetiveType) error
	RemoveChatSensetive(chatId int64, sensetiveWord string, sensetiveType sensetivetypes.SensetiveType) error
	GetChatSensetiveWords(chatId int64, sensetiveType sensetivetypes.SensetiveType) models.ChatSensetiveWords
	GetAllChatSensetiveWords(sensetiveType sensetivetypes.SensetiveType) map[int64][]string
}
type ChatAdmin interface {
	//AddAdmin(chatId int64, userId int64) error
	//RemoveAdmin(chatId int64, userId int64) error
	GetAdminInfo(chatId int64, userId int64) (models.Admin, bool, error)
}
type Repository struct {
	driverType    string
	Chat          Chat
	Warn          Warn
	WarnHistory   WarnHistory
	MuteHistory   MuteHistory
	KnownUser     KnownUser
	ChatSensetive ChatSensetive
	ChatAdmin     ChatAdmin
}

func NewRepository(driver iDb) (*Repository, error) {
	fmt.Printf("driver type: %T\n", driver.GetDriverImplementation())
	//do not forget pointer to db
	switch dr := driver.GetDriverImplementation().(type) {
	// case *mongolib.IMongo:
	// 	return &Repository{
	// 		driverType: dr.DriverType(),
	// 		User:       user.Init(dr),
	// 	}
	case *sqllib.ISql:
		knownUserScheme, err := sqllibKnownUser.Init(dr)
		if err != nil {
			return &Repository{}, err
		}
		ChatScheme, err := sqllibChat.Init(dr)
		if err != nil {
			return &Repository{}, err
		}
		WarnScheme, err := sqllibWarn.Init(dr)
		if err != nil {
			return &Repository{}, err
		}
		WarnHistoryScheme, err := sqllibWarnHistory.Init(dr)
		if err != nil {
			return &Repository{}, err
		}
		MuteHistoryScheme, err := sqllibMuteHistory.Init(dr)
		if err != nil {
			return &Repository{}, err
		}
		chatSensetivesCheme, err := sqllibChatsensetive.Init(dr)
		if err != nil {
			return &Repository{}, err
		}
		return &Repository{
			driverType:    dr.DriverType(),
			KnownUser:     knownUserScheme,
			Chat:          ChatScheme,
			Warn:          WarnScheme,
			WarnHistory:   WarnHistoryScheme,
			MuteHistory:   MuteHistoryScheme,
			ChatSensetive: chatSensetivesCheme,
		}, nil
	}

	return &Repository{}, fmt.Errorf("db driver type not supported")
}
func (r *Repository) RepositoryType() string {
	return r.driverType
}
