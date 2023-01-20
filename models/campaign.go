package models

import (
	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	Name          string `gorm:"size:50" json:"name" binding:"required"`
	EmailTemplate string `gorm:"type:text" binding:"required"`
}
