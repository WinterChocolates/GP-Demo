package models

type Role struct {
	RoleId      int    `gorm:"primaryKey;column:role_id" json:"role_id"`
	RoleName    string `gorm:"unique;size:50;not null" json:"role_name"`
	Description string `gorm:"size:200" json:"description,omitempty"`

	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"-"`
}
