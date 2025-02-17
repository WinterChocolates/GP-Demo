package models

type Permission struct {
	PermID      uint   `gorm:"primaryKey;column:perm_id" json:"perm_id"`
	PermCode    string `gorm:"unique;size:50;not null" json:"perm_code"`
	Description string `gorm:"size:200" json:"description,omitempty"`
}
