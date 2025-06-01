package models

import "time"

type Password struct {
	Hash      string    `bson:"hash" json:"-"`
	OTP       string    `bson:"otp,omitempty" json:"otp,omitempty"`
	Reset     bool      `bson:"reset" json:"reset"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// NewPassword creates an Password instance with default values
func NewPassword(value string) Password {
	return Password{
		Hash:      value,
		OTP:       "",
		Reset:     false, // Default value
		UpdatedAt: time.Now(),
	}
}
