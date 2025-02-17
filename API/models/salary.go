package models

import "time"

type Salary struct {
	SalaryID    uint      `gorm:"primaryKey;column:salary_id" json:"salary_id"`
	UserID      uint      `gorm:"not null;uniqueIndex:uniq_user_month" json:"user_id"`
	Month       string    `gorm:"size:6;not null;uniqueIndex:uniq_user_month" json:"month"`
	Base        float64   `gorm:"type:decimal(10,2);not null" json:"base"`
	Bonus       float64   `gorm:"type:decimal(10,2);default:0.00" json:"bonus"`
	Deductions  float64   `gorm:"type:decimal(10,2);default:0.00" json:"deductions"`
	PaymentDate time.Time `json:"payment_date,omitempty"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
