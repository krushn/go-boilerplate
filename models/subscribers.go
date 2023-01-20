package models

import (
	"gorm.io/gorm"
)

type Subscribers struct {
	gorm.Model
	Name  string `gorm:"size:50" json:"name" binding:"required"`
	Email string `gorm:"size:100" json:"email" binding:"required"`
}
