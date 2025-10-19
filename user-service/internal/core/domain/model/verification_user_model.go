package model

import "time"

type VerificationUser struct {
	ID        int64 `grom:"primaryKey"`
	UserID    int64 `grom:"index"`
	Token     string
	TokenType string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	User      User `gorm:"foreignKey:UserID"`
}
