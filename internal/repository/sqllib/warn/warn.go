package warn

import "github.com/ponyCorp/rebornPony/internal/repository/sqllib"

type WarnScheme struct {
	Id           int      `gorm:"column:id"`
	UserId       int64    `gorm:"column:user_id"`
	ChatId       int64    `gorm:"column:chat_id"`
	Count        int      `gorm:"column:count"`
	LastWarnTime int64    `gorm:"column:last_warn_time"`
	History      []string `gorm:"foreignKey:WarnID"`
}

type Warn struct {
	driver *sqllib.ISql
}

func Init(driver *sqllib.ISql) (*Warn, error) {
	driver.Driver.AutoMigrate(&WarnScheme{})

	return &Warn{
		driver: driver,
	}, nil
}
