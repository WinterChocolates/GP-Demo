package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string     `gorm:"size:50;uniqueIndex;not null;comment:用户名"`
	Email        string     `gorm:"size:50;uniqueIndex;not null;comment:邮箱"`
	Phone        string     `gorm:"size:20;uniqueIndex;not null;comment:手机号"`
	PasswordHash string     `gorm:"size:60;not null;comment:密码哈希"`
	Usertype     string     `gorm:"type:ENUM('admin','employee','candidate');default:'candidate';index;comment:用户类型"`
	Department   string     `gorm:"size:50;index;comment:所属部门"`
	Position     string     `gorm:"size:50;index;comment:职位"`
	HireDate     *time.Time `gorm:"comment:入职日期"`
	SalaryBase   float64    `gorm:"type:decimal(12,2);comment:基本工资"`
	Active       bool       `gorm:"default:true;index;comment:账户状态"`

	Applications    []Application    `gorm:"foreignKey:UserID"`
	Attendances     []Attendance     `gorm:"foreignKey:UserID"`
	Salaries        []Salary         `gorm:"foreignKey:UserID"`
	TrainingRecords []TrainingRecord `gorm:"foreignKey:UserID"`
	Roles           []Role           `gorm:"many2many:user_roles;"`
}
