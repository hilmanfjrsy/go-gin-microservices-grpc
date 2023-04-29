package models

import (
	"time"

	"gorm.io/gorm"
)

type Products struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"size:255;" json:"description"`
	Price       float64        `gorm:"size:255;not null" json:"price"`
	UserId      uint64         `gorm:"not null" json:"user_id"`
	Thumbnail   string         `gorm:"size:255;" json:"thumbnail"`
	Version     uint64         `gorm:"default:0;not null" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Pagination struct {
	TotalPage   uint64 `json:"total_page"`
	PrevPage    uint64 `json:"prev_page"`
	NextPage    uint64 `json:"next_page"`
	CurrentPage uint64 `json:"current_page"`
}
