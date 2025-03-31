package models

type Email struct {
	Value    string `bson:"value" json:"value" validate:"required,email"`
	OTP      string `bson:"otp,omitempty" json:"otp,omitempty"`
	Verified bool   `bson:"verified" json:"verified"`
}

// NewEmail creates an Email instance with default values
func NewEmail(value string) Email {
	return Email{
		Value:    value,
		OTP:      "",
		Verified: false, // Default value
	}
}