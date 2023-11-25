package sqllib

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ISql struct {
	Driver *gorm.DB
	dbType string
}

const (
	DRIVER_TYPE = "sql"
)

func NewDriver(dbType string) *ISql {
	return &ISql{
		dbType: dbType,
	}
}
func (c *ISql) Connect(path string, dbName string) error {
	fmt.Println("sql connect")
	var typeBD gorm.Dialector
	switch c.dbType {
	case "postgres":
		typeBD = postgres.Open(path)
	case "sqlite":
		typeBD = sqlite.Open(path)
	default:
		return fmt.Errorf("database type not supported or config is broken")
	}

	db, err := gorm.Open(typeBD, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	c.Driver = db
	return nil
}
func (c *ISql) DriverType() string {
	return DRIVER_TYPE
}
func (c *ISql) GetDriverImplementation() interface{} {
	return c
}

func (c *ISql) Disconnect() error {
	fmt.Println("sql disconnect")
	instanceDB, err := c.Driver.DB()
	if err != nil {
		return err
	}

	return instanceDB.Close()
}
