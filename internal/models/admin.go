package models

import "github.com/ponyCorp/rebornPony/internal/models/admintypes"

type Admin struct {
	UserId int64
	ChatId int64
	Level  admintypes.AdminType
	Role   string
}
