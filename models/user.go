package models

import (
	"time"
)

type User struct {
	ID              int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name            string     `gorm:"not null" json:"name" form:"name"`
	Email           string     `gorm:"not null unique" json:"email" form:"email"`
	Password        string     `gorm:"not null" json:"password" form:"password"`
	ConfirmPassowrd string     `sql:"-" json:"confirmPassword" form:"confirmPassword"`
	Admin           bool       `gorm:"not null; default: false" json:"admin" form:"admin"`
	CreatedAt       *time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt" form:"updatedAt"`
}

func (user User) MissingFields() string {
	if user.Email == "" {
		return "email"
	} else if user.Password == "" {
		return "password"
	} else if user.ConfirmPassowrd == "" {
		return "confirmação de senha"
	} else if user.Name == "" {
		return "name"
	}
	return ""
}
