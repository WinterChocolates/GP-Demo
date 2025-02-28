package models

import "gorm.io/gorm"

// Role 角色模型
type Role struct {
	gorm.Model
	Name        string `gorm:"size:50;uniqueIndex;not null;comment:角色名称"`
	Description string `gorm:"size:200;comment:角色描述"`

	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Users       []User       `gorm:"many2many:user_roles;"`
}
