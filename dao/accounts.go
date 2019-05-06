package dao

import (
	"github.com/jinzhu/gorm"
)

// 帐户验证模式定义
const (
	// AVSPassword 密
	AVSPassword string = "password"

	// AVSPhoneCode 手机验证
	AVSPhoneCode = "phone"

	// AVSEmail 邮件验证码
	AVSEmail = "email"
)

type Accounts struct {
	gorm.Model
	Schemes string `gorm:"type:varchar(50)"`                    // 帐户验证模式
	Name    string `gorm:"type:varchar(255);default(password)"` // 名称
	Secret  string `gorm:"type:varchar(255)"`                   // 验证数据
	Status  int8   // 状态
}

func (Accounts) TableName() string {
	return "accounts"
}
