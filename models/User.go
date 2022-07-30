package models

import "time"

type User struct {
	BaseModel
	Username    string    `gorm:"uniqueIndex" json:"username"`
	Password    string    `json:"-"`
	Token       string    `json:"token"`
	TokenExpiry time.Time `json:"token_expiry"`
}
