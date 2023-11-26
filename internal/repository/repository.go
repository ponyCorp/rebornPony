package repository

import (
	"fmt"

	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
	sqllibChat "github.com/ponyCorp/rebornPony/internal/repository/sqllib/chat"
	sqllibUser "github.com/ponyCorp/rebornPony/internal/repository/sqllib/user"
)

type User interface {
	GetUserByID(id string) (models.User, error)
}
type Warn interface {
	IncreaseWarn(userId int64, chatId int64) error
	DecreaseWarn(userId int64, chatId int64) error
	ResetWarn(userId int64, chatId int64) error
	GetWarnByChatId(userId int64, chatId int64) int
	GetWarns(userId int64) []models.Warn
}
type MuteHistory interface {
	AddUserMuteHistory(mute models.MuteHistory) error
	GetUserMuteHistory(userId int64) []models.MuteHistory
	GetUserMuteHistoryByChatId(userId int64, chatId int64) []models.MuteHistory
}
type WarnHistory interface {
	AddUserWarnHistory(warn models.WarnHistory) error
	GetUserWarnHistory(userId int64) []models.WarnHistory
	GetUserWarnHistoryByChatId(userId int64, chatId int64) []models.WarnHistory
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
	driverType string
	User       User
	Chat       Chat
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
		userScheme, err := sqllibUser.Init(dr)
		if err != nil {
			return &Repository{}, err
		}
		ChatScheme, err := sqllibChat.Init(dr)
		if err != nil {
			return &Repository{}, err
		}
		return &Repository{
			driverType: dr.DriverType(),
			User:       userScheme,
			Chat:       ChatScheme,
		}, nil
	}

	return &Repository{}, fmt.Errorf("db driver type not supported")
}
func (r *Repository) RepositoryType() string {
	return r.driverType
}
