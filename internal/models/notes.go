package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model        // Pakai ini agar otomatis punya ID, CreatedAt, UpdatedAt, DeletedAt
	Title      string `gorm:"type:varchar(100);not null" json:"title"`
	Content    string `gorm:"type:text" json:"content"`
	UserID     uint   `gorm:"not null;index" json:"user_id"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
