package repository

import (
	"fmt"

	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

//type dbConstraint   = mongo.IMongo | sql.ISql

type iDb interface {
	Connect(url string, dbName string) error
	Disconnect() error
	DriverType() string
	GetDriverImplementation() interface{}
}

func CreateDriver(dType string) (iDb, error) {
	switch dType {
	// case "mongo":
	// 	return mongolib.NewDriver()
	case "sqlite":
		return sqllib.NewDriver(dType), nil
	}
	return nil, fmt.Errorf("driver type not supported or config is broken")
}

// func Main() {
// 	dr := CreateDriver("mongo")
// 	dr.Connect("mongodb://localhost:27017", "test")

// 	rep := NewRepository(dr)
// 	if rep == nil {
// 		panic("error: repository is nil")
// 	}
// 	fmt.Println(rep.RepositoryType())

// 	dr.Disconnect()
// }
