package repository

import (
	"fmt"

	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
	sqllibChat "github.com/ponyCorp/rebornPony/internal/repository/sqllib/chat"
	sqllibKnownUser "github.com/ponyCorp/rebornPony/internal/repository/sqllib/knownUser"
	sqllibMuteHistory "github.com/ponyCorp/rebornPony/internal/repository/sqllib/muteHistory"

	sqllibWarn "github.com/ponyCorp/rebornPony/internal/repository/sqllib/warn"
	sqllibWarnHistory "github.com/ponyCorp/rebornPony/internal/repository/sqllib/warnHistory"
)

type KnownUser interface {
	AddKnownUser(knownUser models.KnownUser) (models.KnownUser, error)
	RemoveKnownUser(knownUser models.KnownUser) error
	GetKnownUser(userId int64, chatId int64) (models.KnownUser, error)
}
type Warn interface {
	IncreaseWarn(userId int64, chatId int64) error
	DecreaseWarn(userId int64, chatId int64) error
	ResetWarn(userId int64, chatId int64) error
	GetWarnByChatId(userId int64, chatId int64) int
	GetWarns(userId int64) []models.Warn
}
type WarnHistory interface {
	AddUserWarnHistory(warn models.WarnHistory) error
	GetUserWarnHistory(userId int64) []models.WarnHistory
	GetUserWarnHistoryByChatId(userId int64, chatId int64) []models.WarnHistory
}
type MuteHistory interface {
	AddUserMuteHistory(mute models.MuteHistory) error
	GetUserMuteHistory(userId int64) []models.MuteHistory
	GetUserMuteHistoryByChatId(userId int64, chatId int64) []models.MuteHistory
}

type Chat interface {
	GetChatByID(id string) (models.Chat, error)
	GetDisabledChats() []models.Chat
	GetChats() []models.Chat
	GetChildrenChats(parentChatID string) []models.Chat
	UnregisterChat(chat models.Chat)
	RegisterChat(chat models.Chat)
	UpdateChat(chat models.Chat)
	AddChildrenChat(parentChatID string, childrenChatID string)
	SetParentChat(parentChatID string, childrenChatID string)
}
type Repository struct {
	driverType  string
	Chat        Chat
	Warn        Warn
	WarnHistory WarnHistory
	MuteHistory MuteHistory
	KnownUser   KnownUser
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

		return &Repository{
			driverType:  dr.DriverType(),
			KnownUser:   knownUserScheme,
			Chat:        ChatScheme,
			Warn:        WarnScheme,
			WarnHistory: WarnHistoryScheme,
			MuteHistory: MuteHistoryScheme,
		}, nil
	}

	return &Repository{}, fmt.Errorf("db driver type not supported")
}
func (r *Repository) RepositoryType() string {
	return r.driverType
}
