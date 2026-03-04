package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Name       string         `gorm:"not null" json:"name"`
	Age        int            `json:"age"`
	Gender     string         `json:"gender"`
	Department string         `json:"department"`
	Position   string         `json:"position"`
	HireDate   string         `gorm:"type:date" json:"hire_date"` // 改为字符串类型
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Employee) TableName() string {
	return "employees"
}
