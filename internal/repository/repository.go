package repository

import (
	"fmt"

	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
	sqllibUser "github.com/ponyCorp/rebornPony/internal/repository/sqllib/user"
)

type User interface {
	GetUserByID(id string) (models.User, error)
}

type Repository struct {
	driverType string
	User       User
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
		return &Repository{
			driverType: dr.DriverType(),
			User:       userScheme,
		}, nil
	}

	return &Repository{}, fmt.Errorf("db driver type not supported")
}
func (r *Repository) RepositoryType() string {
	return r.driverType
}
