package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"size:255"`
	Username  string `gorm:"size:255;uniqueIndex;not null"`
	Password  string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
