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

type Chat interface {
	GetChatByID(id string) (models.Chat, error)
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
