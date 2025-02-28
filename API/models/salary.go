package models

import (
	"time"

	"gorm.io/gorm"
)

// Salary 薪资模型
type Salary struct {
	gorm.Model
	UserID      uint      `gorm:"uniqueIndex:uniq_user_month;not null;comment:用户ID"`
	Month       string    `gorm:"size:6;uniqueIndex:uniq_user_month;comment:薪资月份"`
	Base        float64   `gorm:"type:decimal(12,2);not null;comment:基本工资"`
	Bonus       float64   `gorm:"type:decimal(12,2);default:0.00;comment:奖金"`
	Deductions  float64   `gorm:"type:decimal(12,2);default:0.00;comment:扣款"`
	PaymentDate time.Time `gorm:"comment:发放日期"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
