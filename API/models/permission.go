package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	PermCode    string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string `gorm:"type:varchar(200)"`

	Roles []Role `gorm:"many2many:role_permissions;"`
}
