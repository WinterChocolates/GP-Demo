package models

import "gorm.io/gorm"

// Permission 权限模型
type Permission struct {
	gorm.Model
	Code        string `gorm:"size:50;uniqueIndex;not null;comment:权限代码"`
	Description string `gorm:"size:200;comment:权限描述"`

	Roles []Role `gorm:"many2many:role_permissions;"`
}
