package model

import (
	"time"
)

// User represents user data
type User struct {
	Acct      string     `gorm:"primary_key"  json:"acct"`
	Pwd       string     `json:"-"`
	Fullname  string     `json:"fullname,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
