package models

import (
	"time"

	"gorm.io/gorm"
)

type Salary struct {
	gorm.Model
	UserID      uint    `gorm:"uniqueIndex:uniq_user_month;not null"`
	Month       string  `gorm:"type:char(6);uniqueIndex:uniq_user_month;comment:'格式: YYYYMM'"`
	Base        float64 `gorm:"type:decimal(10,2);not null"`
	Bonus       float64 `gorm:"type:decimal(10,2);default:0.00"`
	Deductions  float64 `gorm:"type:decimal(10,2);default:0.00"`
	PaymentDate *time.Time

	User User `gorm:"foreignKey:UserID"`
}
