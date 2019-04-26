package dao

import (
	"github.com/jinzhu/gorm"
)

type Accounts struct {
	gorm.Model
	ID uint64 `gorm:AUTO_INCREMENT`
	Username string
	Password string
}