package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	RoleName    string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string `gorm:"type:varchar(200)"`

	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Users       []User       `gorm:"many2many:user_roles;"`
}
