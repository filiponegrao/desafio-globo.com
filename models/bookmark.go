package models

import "time"

type Bookmark struct {
	ID        int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Owner     User       `json:"owner" form:"owner"`
	UserID    int64      `json:"user_id" form:"user_id"`
	URL       string     `json:"url" form:"url"`
	Comment   string     `json:"comment" form:"comment"`
	CreatedAt *time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" form:"updatedAt"`
}
