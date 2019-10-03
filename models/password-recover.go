package models

import "time"

type PasswordRecover struct {
	ID        int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserID    int64      `json:"user_id" form:"user_id"`
	Hash      string     `json:"hash" form:"hash"`
	CreatedAt *time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" form:"updatedAt"`
}
