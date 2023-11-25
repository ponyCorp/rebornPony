package repository

import (
	"fmt"

	"github.com/ponyCorp/rebornPony/internal/repository/mongolib/user"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

type User interface {
	GetName() string
	GetUserByID(id string) (*user.User, error)
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

		return &Repository{
			driverType: dr.DriverType(),
		}, nil
	}

	return &Repository{}, fmt.Errorf("db driver type not supported")
}
func (r *Repository) RepositoryType() string {
	return r.driverType
}
