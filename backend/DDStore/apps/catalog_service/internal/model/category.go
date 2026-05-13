package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	Id        int            `json:"id" gorm:"primaryKey;autoIncrement;<-:create"`
	Name      string         `json:"name" gorm:"column:name"`
	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `json:"delete_at" gorm:"index"`
}
