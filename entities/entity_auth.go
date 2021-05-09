package entities

import "time"

type ExpireType uint8

const (
	Expired ExpireType = iota
	NotExpired
	All
)

type ResetPwdStatus uint8

const (
	ResetPwdStatusNoActionYet ResetPwdStatus = iota
	ResetPwdStatusAccepted
	ResetPwdStatusCanceled
)

type AuthUser struct {
	User
	Password []byte
}

type UserToken struct {
	ID        string
	UserID    string
	Token     string
	Device    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type RefreshToken struct {
	ID        string    `json:"-"`
	UserID    string    `json:"-"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"-"`
}

type ResetPwdReq struct {
	ID           string         `json:"id"`
	EmailAddress string         `json:"emailAddress"`
	Status       ResetPwdStatus `json:"status"`
	ExpiresAt    time.Time      `json:"expiresAt"`
	CreatedAt    time.Time      `json:"createdAt"`
}
