package user

import (
	"github.com/ponyCorp/rebornPony/internal/models"
	"github.com/ponyCorp/rebornPony/internal/repository/sqllib"
)

type User struct {
	driver *sqllib.ISql
}
type UserScheme struct {
	ID            int64  `gorm:"column:id"`
	UserID        int64  `gorm:"column:user_id"`
	ChatID        int64  `gorm:"column:chat_id"`
	Name          string `gorm:"column:name"`
	Warn          int    `gorm:"column:warn"`
	MutesFromWarn int    `gorm:"column:mutes_from_warn"`
}

func Init(driver *sqllib.ISql) (*User, error) {
	driver.Driver.AutoMigrate(&UserScheme{})
	return &User{
		driver: driver,
	}, nil
}

// GetUserByID
func (u *User) GetUserByID(userID string) (models.User, error) {
	var user UserScheme
	if result := u.driver.Driver.Where("user_id = ?", userID).First(&user); result.Error != nil {

		// if result.Error.Error() == "record not found" {
		// 	u.driver.Driver.Where("user_id = ?", 0).First(&user)
		// }
		return models.User{}, result.Error
	}

	return models.User{
		Name: user.Name,
	}, nil

}
