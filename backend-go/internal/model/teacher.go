package model

import "time"

type Teacher struct {
	ID                  int        `json:"id"`
	Email               string     `json:"email"`
	EncryptedPassword   string     `json:"-"`
	ResetPasswordToken  *string    `json:"-"`
	ResetPasswordSentAt *time.Time `json:"-"`
	RememberCreatedAt   *time.Time `json:"-"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	Name                *string    `json:"name"`
}