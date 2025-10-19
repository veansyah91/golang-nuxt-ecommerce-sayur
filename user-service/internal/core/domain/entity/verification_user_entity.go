package entity

import "time"

type VerificationUserEntity struct {
	ID        int64
	UserID    int64
	Token     string
	TokenType string
	ExpiresAt time.Time
	User      UserEntity
}
